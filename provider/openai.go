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

	prompt, err := p.buildPrompt(msg)
	if err != nil {
		return 0, err
	}

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
			},
		},
	)

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
