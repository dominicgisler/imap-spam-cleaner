package provider

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/ollama/ollama/api"
)

type Ollama struct {
	AIBase
	client *api.Client
	url    *url.URL
}

func (p *Ollama) Name() string {
	return "ollama"
}

func (p *Ollama) ValidateConfig(config map[string]string) error {

	if err := p.AIBase.ValidateConfig(config); err != nil {
		return err
	}

	if config["url"] == "" {
		return errors.New("ollama url is required")
	}

	u, err := url.Parse(config["url"])
	if err != nil {
		return err
	}
	p.url = u

	return nil
}

func (p *Ollama) Init(config map[string]string) error {
	if err := p.ValidateConfig(config); err != nil {
		return err
	}
	p.client = api.NewClient(p.url, http.DefaultClient)
	return nil
}

func (p *Ollama) Analyze(msg imap.Message) (int, error) {

	userContent, err := p.buildUserPrompt(msg)
	if err != nil {
		return 0, err
	}

	messages := []api.Message{}
	if p.systemPrompt != "" {
		messages = append(messages, api.Message{
			Role:    "system",
			Content: p.systemPrompt,
		})
	}
	messages = append(messages, api.Message{
		Role:    "user",
		Content: userContent,
	})

	b := false
	req := api.ChatRequest{
		Model:    p.model,
		Messages: messages,
		Stream:   &b,
	}

	if p.temperature != nil || p.topP != nil || p.maxTokens != nil {
		opts := map[string]any{}
		if p.temperature != nil {
			opts["temperature"] = *p.temperature
		}
		if p.topP != nil {
			opts["top_p"] = *p.topP
		}
		if p.maxTokens != nil {
			opts["num_predict"] = *p.maxTokens
		}
		req.Options = opts
	}

	var resp string
	if err = p.client.Chat(context.Background(), &req, func(response api.ChatResponse) error {
		resp = response.Message.Content
		return nil
	}); err != nil {
		return 0, err
	}

	i, err := strconv.ParseInt(resp, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}
