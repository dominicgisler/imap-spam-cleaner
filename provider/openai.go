package provider

import (
	"context"
	"errors"
	"strconv"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	AIBase
	client *openai.Client
	apikey string
}

func (p *OpenAI) Name() string {
	return "openai"
}

func (p *OpenAI) ValidateConfig(config map[string]string) error {

	if err := p.AIBase.ValidateConfig(config); err != nil {
		return err
	}

	if config["apikey"] == "" {
		return errors.New("openai apikey is required")
	}
	p.apikey = config["apikey"]

	return nil
}

func (p *OpenAI) Init(config map[string]string) error {
	if err := p.ValidateConfig(config); err != nil {
		return err
	}
	p.client = openai.NewClient(p.apikey)
	return nil
}

func (p *OpenAI) Analyze(msg imap.Message) (int, error) {

	userContent, err := p.buildUserPrompt(msg)
	if err != nil {
		return 0, err
	}

	messages := []openai.ChatCompletionMessage{}
	if p.systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: p.systemPrompt,
		})
	}
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userContent,
	})

	req := openai.ChatCompletionRequest{
		Model:    p.model,
		Messages: messages,
	}
	if p.temperature != nil {
		req.Temperature = *p.temperature
	}
	if p.topP != nil {
		req.TopP = *p.topP
	}
	if p.maxTokens != nil {
		req.MaxCompletionTokens = int(*p.maxTokens)
	}

	resp, err := p.client.CreateChatCompletion(context.Background(), req)

	if err != nil {
		return 0, err
	}

	if len(resp.Choices) == 0 {
		return 0, errors.New("empty openai response")
	}

	i, err := strconv.ParseInt(resp.Choices[0].Message.Content, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}
