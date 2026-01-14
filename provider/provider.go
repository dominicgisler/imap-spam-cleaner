package provider

import (
	"errors"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
)

type Provider interface {
	Name() string
	Init(config map[string]string) error
	ValidateConfig(config map[string]string) error
	Analyze(message imap.Message) (int, error)
}

func New(t string) (Provider, error) {
	providers := []Provider{&OpenAI{}, &Ollama{}, &SpamAssassin{}}
	for _, provider := range providers {
		if provider.Name() == t {
			return provider, nil
		}
	}
	return nil, errors.New("unknown provider")
}
