package config

import (
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Providers map[string]Provider `yaml:"providers" validate:"required,dive"`
	Inboxes   []Inbox             `yaml:"inboxes"   validate:"required,dive"`
}

type Provider struct {
	Type        string            `yaml:"type"        validate:"required,oneof=openai"`
	Credentials map[string]string `json:"credentials" validate:"required"`
}

type Inbox struct {
	Host     string `yaml:"host"     validate:"required"`
	Port     int    `yaml:"port"     validate:"required"`
	TLS      bool   `yaml:"tls"      validate:"omitempty"`
	Username string `yaml:"username" validate:"required"`
	Password string `yaml:"password" validate:"required"`
	Provider string `yaml:"provider" validate:"required"`
	Inbox    string `yaml:"inbox"    validate:"required"`
	Spam     string `yaml:"spam"     validate:"required"`
	MinScore int    `yaml:"minscore" validate:"required"`
}

func Load(path string) (*Config, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	if err = validator.New().Struct(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
