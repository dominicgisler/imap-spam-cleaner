package app

import (
	"github.com/dominicgisler/imap-spam-cleaner/config"
)

type Options struct {
	RunNow      bool
	AnalyzeOnly bool
}

type Context struct {
	Config  *config.Config
	Options Options
}
