package provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/sashabaranov/go-openai"
	"strconv"
)

type OpenAI struct {
	client  *openai.Client
	apikey  string
	model   string
	maxsize int
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

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: p.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
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
