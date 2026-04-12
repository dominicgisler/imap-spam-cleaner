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

// textBodyDivisor controls the plain-text share of the LLM prompt budget when
// both text and HTML bodies are present: plain-text gets 1/textBodyDivisor of
// maxsize; the rest goes to the HTML-derived Markdown, reflecting that spam
// signals tend to be denser in the HTML part.
const textBodyDivisor = 4

const defaultSystemPrompt = `You are a spam classification assistant. Analyze emails objectively and return only a single integer score.`

const defaultUserPrompt = `
Analyze the following email for its spam potential.
Return a spam score between 0 and 100. Only answer with the number itself, no other text.

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

type AIBase struct {
	model        string
	maxsize      int
	systemPrompt string
	userPrompt   *template.Template
	temperature  *float32
	topP         *float32
	maxTokens    *int32
}

func (p *AIBase) ValidateConfig(config map[string]string) error {

	if config["model"] == "" {
		return errors.New("ai model is required")
	}
	p.model = config["model"]

	n, err := strconv.Atoi(config["maxsize"])
	if err != nil || n < 1 {
		return errors.New("maxsize must be a positive integer")
	}
	p.maxsize = n

	p.systemPrompt = defaultSystemPrompt
	if config["system_prompt"] != "" {
		p.systemPrompt = config["system_prompt"]
	}

	// "user_prompt" is the canonical key; "prompt" is kept for backward compatibility.
	userPromptStr := defaultUserPrompt
	if config["user_prompt"] != "" {
		userPromptStr = config["user_prompt"]
	} else if config["prompt"] != "" {
		userPromptStr = config["prompt"]
	}

	p.userPrompt, err = template.New("user_prompt").Parse(userPromptStr)
	if err != nil {
		return err
	}

	if s := config["temperature"]; s != "" {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return errors.New("temperature must be a float")
		}
		v := float32(f)
		p.temperature = &v
	}

	if s := config["top_p"]; s != "" {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return errors.New("top_p must be a float")
		}
		v := float32(f)
		p.topP = &v
	}

	if s := config["max_tokens"]; s != "" {
		n, err := strconv.ParseInt(s, 10, 32)
		if err != nil || n < 1 {
			return errors.New("max_tokens must be a positive integer")
		}
		v := int32(n)
		p.maxTokens = &v
	}

	return nil
}

func (p *AIBase) buildUserPrompt(msg imap.Message) (string, error) {

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
		textLimit := p.maxsize / textBodyDivisor
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
	if err := p.userPrompt.Execute(&buf, TplVars{
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
		return "", errors.New("user_prompt template error: " + err.Error())
	}

	return buf.String(), nil
}
