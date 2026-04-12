package inbox

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dominicgisler/imap-spam-cleaner/app"
	"github.com/dominicgisler/imap-spam-cleaner/checkpoint"
	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/dominicgisler/imap-spam-cleaner/provider"
	goimap "github.com/emersion/go-imap/v2"
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

func processInbox(ctx app.Context, inboxCfg app.Inbox, prov app.Provider) {

	var err error
	var p provider.Provider
	var im *imap.Imap

	logx.Infof("Handling %s", inboxCfg.Username)

	cp, err := checkpoint.Load(inboxCfg.Host, inboxCfg.Username, inboxCfg.Inbox)
	if err != nil {
		logx.Errorf("Could not load checkpoint: %v\n", err)
		return
	}

	if im, err = imap.New(inboxCfg); err != nil {
		logx.Errorf("Could not load imap: %v\n", err)
		return
	}
	defer im.Close()

	currentUIDValidity := im.GetUIDValidity()

	// First run: no checkpoint exists yet — establish the baseline UID.
	if cp == nil {
		logx.Infof("First run for %s: establishing checkpoint (no messages processed this run)", inboxCfg.Username)
		maxUID, err := im.GetMaxUID()
		if err != nil {
			logx.Errorf("Could not get max UID: %v\n", err)
			return
		}
		if err = checkpoint.Save(inboxCfg.Host, inboxCfg.Username, inboxCfg.Inbox, &checkpoint.Checkpoint{
			UIDValidity: currentUIDValidity,
			LastUID:     uint32(maxUID),
		}); err != nil {
			logx.Errorf("Could not save checkpoint: %v\n", err)
		}
		logx.Infof("Checkpoint initialised at UID %d (UIDValidity=%d)", maxUID, currentUIDValidity)
		return
	}

	// UIDVALIDITY changed: UIDs are no longer meaningful — reset the checkpoint.
	if cp.UIDValidity != currentUIDValidity {
		logx.Warnf("UIDVALIDITY changed for %s (%d → %d): resetting checkpoint", inboxCfg.Username, cp.UIDValidity, currentUIDValidity)
		maxUID, err := im.GetMaxUID()
		if err != nil {
			logx.Errorf("Could not get max UID after UIDVALIDITY change: %v\n", err)
			return
		}
		if err = checkpoint.Save(inboxCfg.Host, inboxCfg.Username, inboxCfg.Inbox, &checkpoint.Checkpoint{
			UIDValidity: currentUIDValidity,
			LastUID:     uint32(maxUID),
		}); err != nil {
			logx.Errorf("Could not save checkpoint after UIDVALIDITY reset: %v\n", err)
		}
		logx.Infof("Checkpoint reset to UID %d (UIDValidity=%d)", maxUID, currentUIDValidity)
		return
	}

	// Incremental run: process only messages newer than the last checkpoint.
	sinceUID := goimap.UID(cp.LastUID)
	msgs, err := im.LoadMessages(sinceUID)
	if err != nil {
		logx.Errorf("Could not load messages: %v\n", err)
		return
	}
	logx.Infof("Loaded %d new messages since UID %d", len(msgs), sinceUID)

	p, err = provider.New(prov.Type)
	if err != nil {
		logx.Errorf("Could not load provider: %v\n", err)
		return
	}

	if err = p.Init(prov.Config); err != nil {
		logx.Errorf("Could not init provider: %v\n", err)
		return
	}

	moved := 0
	for _, m := range msgs {
		if wl, ok := ctx.Config.Whitelists[inboxCfg.Whitelist]; ok {
			trustedSender := false
			for _, rgx := range wl {
				if rgx.Match([]byte(m.From)) {
					trustedSender = true
					break
				}
			}
			if trustedSender {
				logx.Debugf("Skipping message #%d (%s) because of trusted sender (%s)", m.UID, m.Subject, m.From)
				// Advance checkpoint for skipped (trusted) messages.
				if err = checkpoint.Save(inboxCfg.Host, inboxCfg.Username, inboxCfg.Inbox, &checkpoint.Checkpoint{
					UIDValidity: currentUIDValidity,
					LastUID:     uint32(m.UID),
				}); err != nil {
					logx.Errorf("Could not save checkpoint for UID %d: %v\n", m.UID, err)
				}
				continue
			}
		}

		n, err := p.Analyze(m)
		if err != nil {
			logx.Errorf("Could not analyze message (%s): %v\n", m.Subject, err)
			// Do not advance checkpoint — retry on next run.
			continue
		}
		logx.Debugf("Spam score of message #%d (%s): %d/100", m.UID, m.Subject, n)

		if n >= inboxCfg.MinScore {
			if ctx.Options.AnalyzeOnly {
				logx.Debugf("Analyze only mode, not moving message #%d", m.UID)
			} else {
				if err = im.MoveMessage(m.UID, inboxCfg.Spam); err != nil {
					logx.Errorf("Could not move message (%s): %v\n", m.Subject, err)
					// Do not advance checkpoint — retry on next run.
					continue
				}
				moved++
			}
		}

		// Advance checkpoint only after the message was fully and successfully processed.
		if err = checkpoint.Save(inboxCfg.Host, inboxCfg.Username, inboxCfg.Inbox, &checkpoint.Checkpoint{
			UIDValidity: currentUIDValidity,
			LastUID:     uint32(m.UID),
		}); err != nil {
			logx.Errorf("Could not save checkpoint for UID %d: %v\n", m.UID, err)
		}
	}
	logx.Infof("Moved %d messages", moved)
}
