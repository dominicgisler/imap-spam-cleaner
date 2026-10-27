package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"text/template"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
)

type TypeSafe struct {
	AIBase
	apikey       string
	url          *url.URL
	instructions string
}

type TypeSafeQuestion struct {
	Type         string `json:"type"`
	Instructions string `json:"instructions"`
}

type TypeSafeRequest struct {
	Model     string                      `json:"model"`
	State     string                      `json:"state"`
	Questions map[string]TypeSafeQuestion `json:"questions"`
}

type TypeSafeResponse struct {
	Answers map[string]struct {
		Noul *float64 `json:"noul"`
	} `json:"answers"`
}

func (p *TypeSafe) Name() string {
	return "typesafe"
}

func (p *TypeSafe) ValidateConfig(config map[string]string) error {

	if err := p.AIBase.ValidateConfig(config); err != nil {
		return err
	}

	// the question is sent separately, so the prompt only carries the message
	if config["prompt"] == "" {
		p.prompt = template.Must(template.New("prompt").Parse(`
{{.RawHeader}}

Content:
{{.Content}}
`))
	}

	if config["apikey"] == "" {
		return errors.New("typesafe apikey is required")
	}
	p.apikey = config["apikey"]

	uri := "https://api.typesafe.ai"
	if config["url"] != "" {
		uri = config["url"]
	}
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	p.url = u

	p.instructions = "Is this email spam, phishing or a scam?"
	if config["instructions"] != "" {
		p.instructions = config["instructions"]
	}

	return nil
}

func (p *TypeSafe) Init(config map[string]string) error {
	if err := p.ValidateConfig(config); err != nil {
		return err
	}
	return nil
}

func (p *TypeSafe) Analyze(msg imap.Message) (int, error) {

	state, err := p.buildPrompt(msg)
	if err != nil {
		return 0, err
	}

	body, err := json.Marshal(TypeSafeRequest{
		Model: p.model,
		State: state,
		Questions: map[string]TypeSafeQuestion{
			"spam": {
				Type:         "noul",
				Instructions: p.instructions,
			},
		},
	})
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		p.url.JoinPath("v1", "systemone").String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", "Bearer "+p.apikey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		detail, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return 0, fmt.Errorf("typesafe returned %s: %s", resp.Status, bytes.TrimSpace(detail))
	}

	var result TypeSafeResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	answer, ok := result.Answers["spam"]
	if !ok || answer.Noul == nil {
		return 0, errors.New("empty typesafe response")
	}

	return int(math.Round(*answer.Noul * 100)), nil
}
