package main

import (
	"flag"
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/inbox"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
)

func init() {
	var v bool
	flag.BoolVar(&v, "verbose", false, "Enable debug logging")
	flag.Parse()
	if v {
		logx.Init(true)
	}
}

func main() {

	c, err := config.Load("config.yml")
	if err != nil {
		panic(err)
	}

	inbox.Process(c)
}
