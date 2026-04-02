package provider

import (
	"context"
	"errors"
	"strconv"

	"github.com/dominicgisler/imap-spam-cleaner/imap"

	"google.golang.org/genai"
)

type Gemini struct {
	AIBase
	client *genai.Client
	apikey string
}

func (p *Gemini) Name() string {
	return "gemini"
}

func (p *Gemini) ValidateConfig(config map[string]string) error {

	if err := p.AIBase.ValidateConfig(config); err != nil {
		return err
	}

	if config["apikey"] == "" {
		return errors.New("gemini apikey is required")
	}
	p.apikey = config["apikey"]

	return nil
}

func (p *Gemini) Init(config map[string]string) error {

	if err := p.ValidateConfig(config); err != nil {
		return err
	}

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: p.apikey,
	})
	if err != nil {
		return err
	}

	p.client = client

	return nil
}

func (p *Gemini) Analyze(msg imap.Message) (int, error) {

	prompt, err := p.buildPrompt(msg)
	if err != nil {
		return 0, err
	}

	resp, err := p.client.Models.GenerateContent(
		context.Background(),
		p.model,
		[]*genai.Content{
			{
				Parts: []*genai.Part{
					{
						Text: prompt,
					},
				},
			},
		},
		nil,
	)

	if err != nil {
		return 0, err
	}

	if len(resp.Candidates) == 0 ||
		len(resp.Candidates[0].Content.Parts) == 0 {
		return 0, errors.New("empty gemini response")
	}

	i, err := strconv.ParseInt(resp.Candidates[0].Content.Parts[0].Text, 10, 64)
	if err != nil {
		return 0, err
	}

	return int(i), nil
}
