package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dominicgisler/imap-spam-cleaner/app"
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/inbox"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/dominicgisler/imap-spam-cleaner/provider"
)

const (
	appName    = "imap-spam-cleaner"
	appVersion = "0.5.0"
)

var options app.Options

func init() {
	var v bool
	flag.BoolVar(&v, "version", false, "Show short version")
	flag.BoolVar(&options.RunNow, "now", false, "Run all inboxes once immediately, ignoring schedule")
	flag.BoolVar(&options.AnalyzeOnly, "analyze-only", false, "Only analyze emails, do not move or delete them")
	flag.Parse()
	if v {
		fmt.Printf("%s v%s\n", appName, appVersion)
		os.Exit(0)
	}
}

func main() {
	logx.Infof("Starting %s v%s", appName, appVersion)
	c, err := config.Load()
	if err != nil {
		logx.Errorf("Could not load config: %v", err)
		return
	}

	ctx := app.Context{
		Config:  c,
		Options: options,
	}

	var p provider.Provider
	for name, prov := range c.Providers {
		p, err = provider.New(prov.Type)
		if err != nil {
			logx.Errorf("Could not load provider: %v\n", err)
			return
		}
		if err = p.ValidateConfig(prov.Config); err != nil {
			logx.Errorf("Invalid config for provider %s: %v\n", name, err)
			return
		}
	}

	if ctx.Options.RunNow {
		logx.Info("Running all inboxes once immediately")
		inbox.RunAllInboxes(ctx)
		return
	}

	inbox.Schedule(ctx)
}
