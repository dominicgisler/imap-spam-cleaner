package config

import (
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
	"os"
)

const configPath = "config.yml"

type Config struct {
	Logging   Logging             `yaml:"logging" validate:"required"`
	Providers map[string]Provider `yaml:"providers" validate:"required,dive"`
	Inboxes   []Inbox             `yaml:"inboxes"   validate:"required,dive"`
}

type Logging struct {
	Level string `yaml:"level" validate:"omitempty"`
}

type Provider struct {
	Type        string            `yaml:"type"        validate:"required,oneof=openai"`
	Credentials map[string]string `yaml:"credentials" validate:"required"`
}

type Inbox struct {
	Schedule string `yaml:"schedule" validate:"required"`
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

func Load() (*Config, error) {

	data, err := os.ReadFile(configPath)
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

	if config.Logging.Level != "" {
		logx.SetLevel(config.Logging.Level)
	}

	logx.Debug("Loaded config")

	return &config, nil
}
