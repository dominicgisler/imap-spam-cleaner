package main

import (
	"flag"
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/inbox"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/sirupsen/logrus"
)

func init() {
	var v bool
	flag.BoolVar(&v, "verbose", false, "Enable debug logging")
	flag.Parse()
	if v {
		logx.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	c, err := config.Load()
	if err != nil {
		logx.Errorf("Could not load config: %v", err)
		return
	}
	inbox.Schedule(c)
}
