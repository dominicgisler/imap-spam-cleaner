package main

import (
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/provider"
	"log"
)

func main() {

	c, err := config.Load("config.yml")
	if err != nil {
		panic(err)
	}

	var msgs []imap.Message
	var p provider.Provider
	for i, inbox := range c.Inboxes {
		prov, ok := c.Providers[inbox.Provider]
		if !ok {
			log.Printf("invalid provider %s for inbox %d", inbox.Provider, i)
			continue
		}
		if msgs, err = imap.LoadMessages(inbox); err != nil {
			log.Printf("could not load messages: %v\n", err)
			continue
		}
		log.Printf("loaded %d messages", len(msgs))
		p, err = provider.New(prov.Type)
		if err != nil {
			log.Printf("could not load provider: %v\n", err)
			continue
		}
		if err = p.Init(prov.Credentials); err != nil {
			log.Printf("could not init provider: %v\n", err)
			continue
		}
		for _, m := range msgs {
			n, err := p.Analyze(m)
			if err != nil {
				log.Printf("could not analyze message: %v\n", err)
			}
			log.Printf("Spam score of \"%s\": %d\n", m.Subject, n)
		}
	}
}
