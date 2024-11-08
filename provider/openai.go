package provider

import (
	"context"
	"fmt"
	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/sashabaranov/go-openai"
	"strconv"
	"strings"
)

type OpenAI struct {
	client *openai.Client
}

func (p *OpenAI) Name() string {
	return "openai"
}

func (p *OpenAI) Init(credentials map[string]string) error {
	p.client = openai.NewClient(credentials["apikey"])
	return nil
}

func (p *OpenAI) Analyze(msg imap.Message) (int, error) {

	resp, err := p.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
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
						strings.Join(msg.Contents, "\n"),
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
