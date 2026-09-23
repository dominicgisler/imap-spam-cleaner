package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/sashabaranov/go-openai"
)

type OpenAI struct {
	AIBase
	client *openai.Client
	apikey string
	url    string
}

func (p *OpenAI) Name() string {
	return "openai"
}

func (p *OpenAI) ValidateConfig(config map[string]string) error {

	if err := p.AIBase.ValidateConfig(config); err != nil {
		return err
	}

	p.url = config["url"]
	p.apikey = config["apikey"]

	if p.url == "" && p.apikey == "" {
		return errors.New("openai apikey is required")
	}
	if p.apikey == "" {
		p.apikey = "lm-studio"
	}

	return nil
}

func (p *OpenAI) Init(config map[string]string) error {
	if err := p.ValidateConfig(config); err != nil {
		return err
	}

	cfg := openai.DefaultConfig(p.apikey)
	if p.url != "" {
		cfg.BaseURL = p.url
	}
	p.client = openai.NewClientWithConfig(cfg)
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
					Role:    openai.ChatMessageRoleUser,
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

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	if idx := strings.LastIndex(content, "</think>"); idx != -1 {
		content = strings.TrimSpace(content[idx+len("</think>"):])
	}

	i, err := strconv.ParseInt(content, 10, 64)
	if err != nil {
		re := regexp.MustCompile(`\b(\d{1,3})\b`)
		matches := re.FindStringSubmatch(content)
		if len(matches) > 1 {
			i, err = strconv.ParseInt(matches[1], 10, 64)
		}
	}
	if err != nil {
		return 0, fmt.Errorf("could not parse spam score from response %q: %w", content, err)
	}

	return int(i), nil
}
