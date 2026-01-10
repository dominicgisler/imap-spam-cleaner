package inbox

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dominicgisler/imap-spam-cleaner/app"
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/dominicgisler/imap-spam-cleaner/provider"
	"github.com/go-co-op/gocron/v2"
)

func Schedule(ctx app.Context) {

	s, err := gocron.NewScheduler()
	if err != nil {
		logx.Errorf("Could not create scheduler: %v", err)
		return
	}

	for i, inbox := range ctx.Config.Inboxes {
		logx.Infof("Scheduling inbox %s (%s)", inbox.Username, inbox.Schedule)
		prov, ok := ctx.Config.Providers[inbox.Provider]
		if !ok {
			logx.Errorf("Invalid provider %s for inbox %d", inbox.Provider, i)
			continue
		}
		if _, err = s.NewJob(
			gocron.CronJob(inbox.Schedule, false),
			gocron.NewTask(processInbox, ctx, inbox, prov),
		); err != nil {
			logx.Errorf("Could not schedule inbox %s (%s): %v", inbox.Username, inbox.Schedule, err)
		}
	}

	s.Start()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	sig := <-ch
	logx.Debugf("Received %s, shutting down", sig.String())

	if err = s.Shutdown(); err != nil {
		logx.Errorf("Could not shutdown scheduler: %v ", err)
	}
}

func RunAllInboxes(ctx app.Context) {
	for i, inbox := range ctx.Config.Inboxes {
		logx.Infof("Processing inbox %s", inbox.Username)
		prov, ok := ctx.Config.Providers[inbox.Provider]
		if !ok {
			logx.Errorf("Invalid provider %s for inbox %d", inbox.Provider, i)
			continue
		}
		processInbox(ctx, inbox, prov)
	}
}

func processInbox(ctx app.Context, inbox config.Inbox, prov config.Provider) {

	var err error
	var msgs []imap.Message
	var p provider.Provider
	var im *imap.Imap
	var n int

	logx.Infof("Handling %s", inbox.Username)

	if im, err = imap.New(inbox); err != nil {
		logx.Errorf("Could not load imap: %v\n", err)
		return
	}

	if msgs, err = im.LoadMessages(); err != nil {
		logx.Errorf("Could not load messages: %v\n", err)
		im.Close()
		return
	}
	logx.Infof("Loaded %d messages", len(msgs))

	p, err = provider.New(prov.Type)
	if err != nil {
		logx.Errorf("Could not load provider: %v\n", err)
		im.Close()
		return
	}

	if err = p.Init(prov.Config); err != nil {
		logx.Errorf("Could not init provider: %v\n", err)
		im.Close()
		return
	}

	moved := 0
	for _, m := range msgs {
		if n, err = p.Analyze(m); err != nil {
			logx.Errorf("Could not analyze message (%s): %v\n", m.Subject, err)
			continue
		}
		logx.Debugf("Spam score of message #%d (%s): %d/100", m.UID, m.Subject, n)

		if n >= inbox.MinScore {
			if ctx.Options.AnalyzeOnly {
				logx.Debugf("Analyze only mode, not moving message #%d", m.UID)
				continue
			}

			if err = im.MoveMessage(m.UID, inbox.Spam); err != nil {
				logx.Errorf("Could not move message (%s): %v\n", m.Subject, err)
				continue
			}
			moved++
		}
	}
	logx.Infof("Moved %d messages", moved)

	im.Close()
}
