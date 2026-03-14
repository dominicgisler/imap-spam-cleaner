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

	if config["apikey"] == "" {
		return errors.New("openai apikey is required")
	}
	p.apikey = config["apikey"]

	if config["model"] == "" {
		return errors.New("openai model is required")
	}
	p.model = config["model"]

	n, err := strconv.ParseInt(config["maxsize"], 10, 64)
	if err != nil || n < 1 {
		return errors.New("openai maxsize must be a positive integer")
	}
	p.maxsize = int(n)

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

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: p.buildPrompt(msg),
				},
			},
		},
	)

	if err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(resp.Choices[0].Message.Content, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}
