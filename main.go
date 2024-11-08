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

	var msgs []imap.Message
	for _, i := range c.Inboxes {
		if msgs, err = imap.LoadMessages(i); err != nil {
			log.Println(err)
			continue
		}
		log.Printf("loaded %d messages", len(msgs))
		for _, m := range msgs {
			log.Println("-", m.Subject)
		}
	}
}
