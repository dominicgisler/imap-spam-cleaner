package main

import (
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"log"
)

func main() {

	c, err := config.Load("config.yml")
	if err != nil {
		panic(err)
	}

	for _, i := range c.Inboxes {
		if err = imap.CheckInbox(i); err != nil {
			log.Println(err)
		}
	}
}
