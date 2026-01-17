package provider

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/ollama/ollama/api"
)

type Ollama struct {
	client  *api.Client
	url     *url.URL
	model   string
	maxsize int
}

func (p *Ollama) Name() string {
	return "ollama"
}

func (p *Ollama) ValidateConfig(config map[string]string) error {

	if config["url"] == "" {
		return errors.New("ollama url is required")
	}

	u, err := url.Parse(config["url"])
	if err != nil {
		return err
	}
	p.url = u

	if config["model"] == "" {
		return errors.New("ollama model is required")
	}
	p.model = config["model"]

	n, err := strconv.ParseInt(config["maxsize"], 10, 64)
	if err != nil || n < 1 {
		return errors.New("ollama maxsize must be a positive integer")
	}
	p.maxsize = int(n)

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

	cont := ""
	contLen := 0
	for _, cnt := range msg.Contents {
		contLen += len(cnt)
		if contLen > p.maxsize {
			logx.Debugf("skipping bytes for message #%d (%s)", msg.UID, msg.Subject)
			break
		}
		cont += cnt + "\n"
	}

	b := false
	req := api.ChatRequest{
		Model: p.model,
		Messages: []api.Message{
			{
				Role: "system",
				Content: fmt.Sprintf(
					"Analyze the following email for its spam potential.\n"+
						"Return a spam score between 0 and 100. Only answer with the number itself.\n\n"+
						"From: %s\nTo: %s\nDelivered-To: %s\nCc: %s\nBcc: %s\nSubject: %s\n\n"+
						"Content:\n%s",
					msg.From,
					msg.To,
					msg.DeliveredTo,
					msg.Cc,
					msg.Bcc,
					msg.Subject,
					cont,
				),
			},
		},
		Stream: &b,
	}

	var resp string
	if err := p.client.Chat(context.Background(), &req, func(response api.ChatResponse) error {
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
