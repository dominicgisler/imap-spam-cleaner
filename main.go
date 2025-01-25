package main

import (
	"flag"
	"fmt"
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/inbox"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"os"
)

const (
	appName    = "imap-spam-cleaner"
	appVersion = "0.1.0"
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
	inbox.Schedule(c)
}
