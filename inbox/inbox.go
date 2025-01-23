package inbox

import (
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/dominicgisler/imap-spam-cleaner/provider"
)

func Process(cfg *config.Config) {

	var err error
	var msgs []imap.Message
	var p provider.Provider
	var im *imap.Imap
	var n int

	for i, inbox := range cfg.Inboxes {
		logx.Infof("handling %s", inbox.Username)

		prov, ok := cfg.Providers[inbox.Provider]
		if !ok {
			logx.Warnf("invalid provider %s for inbox %d", inbox.Provider, i)
			continue
		}

		if im, err = imap.New(inbox); err != nil {
			logx.Errorf("could not load imap: %v\n", err)
			continue
		}

		if msgs, err = im.LoadMessages(); err != nil {
			logx.Errorf("could not load messages: %v\n", err)
			im.Close()
			continue
		}
		logx.Infof("loaded %d messages", len(msgs))

		p, err = provider.New(prov.Type)
		if err != nil {
			logx.Errorf("could not load provider: %v\n", err)
			im.Close()
			continue
		}

		if err = p.Init(prov.Credentials); err != nil {
			logx.Errorf("could not init provider: %v\n", err)
			im.Close()
			continue
		}

		moved := 0
		for _, m := range msgs {
			if n, err = p.Analyze(m); err != nil {
				logx.Errorf("could not analyze message (%s): %v\n", m.Subject, err)
				continue
			}
			logx.Debugf("spam score of message #%d (%s): %d/100", m.UID, m.Subject, n)

			if n >= inbox.MinScore {
				if err = im.MoveMessage(m.UID, inbox.Spam); err != nil {
					logx.Errorf("could not move message (%s): %v\n", m.Subject, err)
					continue
				}
				moved++
			}
		}
		logx.Infof("moved %d messages", moved)

		im.Close()
	}
}
