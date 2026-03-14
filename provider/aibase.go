package provider

import (
	"bytes"
	"errors"
	"strconv"
	"text/template"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
)

type AIBase struct {
	model   string
	maxsize int
	prompt  *template.Template
}

func (p *AIBase) ValidateConfig(config map[string]string) error {

	if config["model"] == "" {
		return errors.New("ai model is required")
	}
	p.model = config["model"]

	n, err := strconv.ParseInt(config["maxsize"], 10, 64)
	if err != nil || n < 1 {
		return errors.New("maxsize must be a positive integer")
	}
	p.maxsize = int(n)

	prompt := `
Analyze the following email for its spam potential.
Return a spam score between 0 and 100. Only answer with the number itself.

From: {{.From}}
To: {{.To}}
Delivered-To: {{.DeliveredTo}}
Cc: {{.Cc}}
Bcc: {{.Bcc}}
Subject: {{.Subject}}

Content:
{{.Content}}
`
	if config["prompt"] != "" {
		prompt = config["prompt"]
	}

	p.prompt, err = template.New("prompt").Parse(prompt)
	if err != nil {
		return err
	}

	return nil
}

func (p *AIBase) buildPrompt(msg imap.Message) (string, error) {

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

	type TplVars struct {
		From        string
		To          string
		DeliveredTo string
		Cc          string
		Bcc         string
		Subject     string
		Content     string
	}

	var buf bytes.Buffer
	if err := p.prompt.Execute(&buf, TplVars{
		From:        msg.From,
		To:          msg.To,
		DeliveredTo: msg.DeliveredTo,
		Cc:          msg.Cc,
		Bcc:         msg.Bcc,
		Subject:     msg.Subject,
		Content:     cont,
	}); err != nil {
		return "", errors.New("prompt template error: " + err.Error())
	}

	return buf.String(), nil
}
