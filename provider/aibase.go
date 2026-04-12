package provider

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"text/template"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/dominicgisler/imap-spam-cleaner/mailclean"
)

// textBodyFraction is the share of the LLM prompt budget allocated to the
// plain-text part of an email when both text and HTML bodies are present.
// The remaining (1 - 1/textBodyFraction) is given to the HTML-derived
// Markdown, reflecting that spam signals tend to be denser in the HTML part.
const textBodyFraction = 4

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

Headers:
{{.Headers}}

From: {{.From}}
To: {{.To}}
Delivered-To: {{.DeliveredTo}}
Cc: {{.Cc}}
Bcc: {{.Bcc}}
Subject: {{.Subject}}

Text body:
{{.TextBody}}

HTML body (converted to Markdown):
{{.HtmlBody}}
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

	textBody := msg.TextBody
	htmlBody := msg.HtmlBody

	// Convert HTML body to simplified Markdown to reduce noise and token count.
	// Falls back to the raw HTML if conversion fails.
	if htmlBody != "" {
		md, err := mailclean.HTMLToSimpleMarkdown(strings.NewReader(htmlBody))
		if err != nil {
			logx.Debugf("HTML to Markdown conversion failed for message #%d (%s), using raw HTML: %v", msg.UID, msg.Subject, err)
		} else {
			htmlBody = md
		}
	}

	// Apply size limits. When both bodies are present, allocate 1/4 of the
	// budget to plain-text and 3/4 to the HTML-derived Markdown — spam
	// signals tend to be denser in the HTML part.
	if textBody != "" && htmlBody != "" {
		textLimit := p.maxsize / textBodyFraction
		htmlLimit := p.maxsize - textLimit
		if len(textBody) > textLimit {
			textBody = textBody[:textLimit]
			logx.Debugf("truncating text body for message #%d (%s)", msg.UID, msg.Subject)
		}
		if len(htmlBody) > htmlLimit {
			htmlBody = htmlBody[:htmlLimit]
			logx.Debugf("truncating HTML body for message #%d (%s)", msg.UID, msg.Subject)
		}
	} else {
		if len(textBody) > p.maxsize {
			textBody = textBody[:p.maxsize]
			logx.Debugf("truncating text body for message #%d (%s)", msg.UID, msg.Subject)
		}
		if len(htmlBody) > p.maxsize {
			htmlBody = htmlBody[:p.maxsize]
			logx.Debugf("truncating HTML body for message #%d (%s)", msg.UID, msg.Subject)
		}
	}

	type TplVars struct {
		From        string
		To          string
		DeliveredTo string
		Cc          string
		Bcc         string
		Subject     string
		Headers     string
		TextBody    string
		HtmlBody    string
	}

	var buf bytes.Buffer
	if err := p.prompt.Execute(&buf, TplVars{
		From:        msg.From,
		To:          msg.To,
		DeliveredTo: msg.DeliveredTo,
		Cc:          msg.Cc,
		Bcc:         msg.Bcc,
		Subject:     msg.Subject,
		Headers:     msg.Headers,
		TextBody:    textBody,
		HtmlBody:    htmlBody,
	}); err != nil {
		return "", errors.New("prompt template error: " + err.Error())
	}

	return buf.String(), nil
}
