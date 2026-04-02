package app

type Options struct {
	RunNow      bool
	AnalyzeOnly bool
}

type Context struct {
	Config  *Config
	Options Options
}
