package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
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

	if config["timeout"] == "" {
		p.timeout = 5 * time.Second
	} else if to, err := time.ParseDuration(config["timeout"]); err == nil && to > 0 {
		p.timeout = to
	} else {
		t, err := strconv.ParseFloat(config["timeout"], 64)
		if err != nil || t <= 0 {
			return errors.New("rspamd timeout must be a duration (eg. 10s, 1m) or a positive number of seconds")
		}
		p.timeout = time.Duration(t * float64(time.Second))
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
