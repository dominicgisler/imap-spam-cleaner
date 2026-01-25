package config

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

const configPath = "config.yml"

type Config struct {
	Logging    Logging                    `yaml:"logging"    validate:"required"`
	Providers  map[string]Provider        `yaml:"providers"  validate:"required,dive"`
	Whitelists map[string][]regexp.Regexp `yaml:"whitelists" validate:"omitempty"`
	Inboxes    []Inbox                    `yaml:"inboxes"    validate:"required,dive"`
}

type Logging struct {
	Level string `yaml:"level" validate:"omitempty"`
}

type Provider struct {
	Type   string            `yaml:"type"   validate:"required,oneof=openai ollama spamassassin"`
	Config map[string]string `yaml:"config" validate:"required"`
}

type Inbox struct {
	Schedule  string        `yaml:"schedule"  validate:"required"`
	Host      string        `yaml:"host"      validate:"required"`
	Port      int           `yaml:"port"      validate:"required"`
	TLS       bool          `yaml:"tls"       validate:"omitempty"`
	Username  string        `yaml:"username"  validate:"required"`
	Password  string        `yaml:"password"  validate:"required"`
	Provider  string        `yaml:"provider"  validate:"required"`
	Inbox     string        `yaml:"inbox"     validate:"required"`
	Spam      string        `yaml:"spam"      validate:"required"`
	MinScore  int           `yaml:"minscore"  validate:"required,gte=0,lte=100"`
	MinAge    time.Duration `yaml:"minage"    validate:"omitempty"`
	MaxAge    time.Duration `yaml:"maxage"    validate:"omitempty"`
	Whitelist string        `yaml:"whitelist" validate:"omitempty"`
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

	for i, inbox := range config.Inboxes {
		if _, ok := config.Providers[inbox.Provider]; !ok {
			return nil, fmt.Errorf("invalid provider %s for inbox #%d", inbox.Provider, i)
		}
		if _, ok := config.Whitelists[inbox.Whitelist]; !ok {
			return nil, fmt.Errorf("invalid whitelist %s for inbox #%d", inbox.Whitelist, i)
		}
	}

	logx.Debug("Loaded config")

	return &config, nil
}
