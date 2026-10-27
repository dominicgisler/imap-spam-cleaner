package provider

import (
	"errors"
	"strconv"
	"time"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
)

type Provider interface {
	Name() string
	Init(config map[string]string) error
	ValidateConfig(config map[string]string) error
	Analyze(message imap.Message) (int, error)
}

func New(t string) (Provider, error) {
	providers := []Provider{&OpenAI{}, &Ollama{}, &SpamAssassin{}, &Gemini{}, &Rspamd{}, &TypeSafe{}}
	for _, provider := range providers {
		if provider.Name() == t {
			return provider, nil
		}
	}
	return nil, errors.New("unknown provider")
}

// parseTimeout accepts a duration (eg. 10s, 1m) or a number of seconds, an empty value returns def
func parseTimeout(value string, def time.Duration) (time.Duration, error) {
	if value == "" {
		return def, nil
	}
	if to, err := time.ParseDuration(value); err == nil && to > 0 {
		return to, nil
	}
	t, err := strconv.ParseFloat(value, 64)
	if err != nil || t <= 0 {
		return 0, errors.New("timeout must be a duration (eg. 10s, 1m) or a positive number of seconds")
	}
	return time.Duration(t * float64(time.Second)), nil
}
