package main

import (
	"flag"
	"fmt"
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/inbox"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/dominicgisler/imap-spam-cleaner/provider"
	"os"
)

const (
	appName    = "imap-spam-cleaner"
	appVersion = "0.3.0"
)

func init() {
	var v bool
	flag.BoolVar(&v, "version", false, "Show short version")
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
	inbox.Schedule(c)
}
