package provider

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

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

func (p *Ollama) HealthCheck(config map[string]string) error {
	if err := p.Init(config); err != nil {
		return err
	}

	host := p.url.Hostname()
	port := p.url.Port()
	if port == "" {
		if p.url.Scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}
	return checkTCP(net.JoinHostPort(host, port), 5*time.Second)
}

func (p *Ollama) Analyze(msg imap.Message) (int, error) {

	prompt, err := p.buildPrompt(msg)
	if err != nil {
		return 0, err
	}

	b := false
	req := api.GenerateRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: &b,
	}

	var resp string
	if err = p.client.Generate(context.Background(), &req, func(response api.GenerateResponse) error {
		resp = response.Response
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
