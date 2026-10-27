package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
)

type Rspamd struct {
	url     string
	timeout time.Duration
}

type RspamdResponse struct {
	Score         float64 `json:"score"`
	RequiredScore float64 `json:"required_score"`
}

func (p *Rspamd) Name() string {
	return "rspamd"
}

func (p *Rspamd) ValidateConfig(config map[string]string) error {

	uri, err := url.Parse(config["url"])
	if err != nil {
		return err
	}
	p.url = uri.String()

	if p.timeout, err = parseTimeout(config["timeout"], 5*time.Second); err != nil {
		return fmt.Errorf("rspamd %w", err)
	}

	return nil
}

func (p *Rspamd) Init(config map[string]string) error {
	if err := p.ValidateConfig(config); err != nil {
		return err
	}
	return nil
}

func (p *Rspamd) Analyze(msg imap.Message) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		p.url+"/checkv2",
		bytes.NewReader(msg.Raw),
	)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "message/rfc822")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("rspamd returned %s", resp.Status)
	}

	var result RspamdResponse
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return int(math.Round(result.Score / result.RequiredScore * 100)), nil
}
