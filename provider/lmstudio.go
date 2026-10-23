package provider

type LMStudio struct {
	OpenAI
}

func (p *LMStudio) Name() string {
	return "lmstudio"
}

func (p *LMStudio) ValidateConfig(config map[string]string) error {
	if config["url"] == "" {
		config["url"] = "http://localhost:1234/v1"
	}
	if config["apikey"] == "" {
		config["apikey"] = "lm-studio"
	}
	return p.OpenAI.ValidateConfig(config)
}
