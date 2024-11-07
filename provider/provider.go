package provider

type Provider interface {
	Init() (*Provider, error)
	Analyze(string) (int, error)
}
