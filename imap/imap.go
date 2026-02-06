package imap

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
)

type Imap struct {
	client *imapclient.Client
	cfg    config.Inbox
}

func New(cfg config.Inbox) (*Imap, error) {

	var err error
	var mailboxes []*imap.ListData

	i := &Imap{
		cfg: cfg,
	}

	if cfg.TLS {
		i.client, err = imapclient.DialTLS(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil)
	} else {
		i.client, err = imapclient.DialInsecure(fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), nil)
	}

	if err != nil {
		i.Close()
		return nil, fmt.Errorf("failed to dial IMAP server: %w", err)
	}

	if err = i.client.Login(cfg.Username, cfg.Password).Wait(); err != nil {
		i.Close()
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	mailboxes, err = i.client.List("", "*", nil).Collect()
	if err != nil {
		return nil, fmt.Errorf("failed to list mailboxes: %w", err)
	}

	logx.Debug("Available mailboxes:")
	for _, l := range mailboxes {
		logx.Debugf("  - %s", l.Mailbox)
	}

	return i, nil
}

func (i *Imap) Close() {
	if i.client != nil {
		i.client.Logout()
		_ = i.client.Close()
	}
}

func (i *Imap) LoadMessages() ([]Message, error) {

	var err error
	var mbox *imap.SelectData
	var msgs []*imapclient.FetchMessageBuffer
	var mr *mail.Reader
	var p *mail.Part
	var messages []Message

	mbox, err = i.client.Select(i.cfg.Inbox, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %w", err)
	}
	logx.Debugf("Found %d messages in inbox", mbox.NumMessages)

	searchCrit := &imap.SearchCriteria{}
	if i.cfg.MinAge > 0 {
		searchCrit.Before = time.Now().Add(-minAge)
	}
	if i.cfg.MaxAge > 0 {
		searchCrit.Since = time.Now().Add(-maxAge)
	}

	uidRes, err := i.client.UIDSearch(searchCrit, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("could not search UIDs: %w", err)
	}

	logx.Debugf("Found %d UIDs in timerange", len(uidRes.AllUIDs()))
	if len(uidRes.AllUIDs()) == 0 {
		return nil, nil
	}

	fetchOptions := &imap.FetchOptions{
		Envelope: true,
		UID:      true,
		BodySection: []*imap.FetchItemBodySection{
			{
				Peek: true,
			},
		},
	}

	msgs, err = i.client.Fetch(uidRes.All, fetchOptions).Collect()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %v", err)
	}

	for _, msg := range msgs {
		var b []byte
		for _, buf := range msg.BodySection {
			b = buf.Bytes
			break
		}

		mr, err = mail.CreateReader(bytes.NewReader(b))
		if err != nil {
			logx.Warnf("failed to create message reader (msg.UID=%d): %v\n", msg.UID, err)
			continue
		}

		message := Message{
			UID:         msg.UID,
			DeliveredTo: mr.Header.Get("Delivered-To"),
			From:        mr.Header.Get("From"),
			To:          mr.Header.Get("To"),
			Cc:          mr.Header.Get("Cc"),
			Bcc:         mr.Header.Get("Bcc"),
			Subject:     msg.Envelope.Subject,
			Contents:    []string{},
			Raw:         b, // Raw original message bytes. Useful for traditional spam filters.
		}

		if message.Date, err = mr.Header.Date(); err != nil {
			logx.Warnf("failed to load message date (msg.UID=%d): %v\n", msg.UID, err)
			continue
		}

		if i.cfg.MinAge > 0 && message.Date.After(time.Now().Add(-i.cfg.MinAge)) || i.cfg.MaxAge > 0 && message.Date.Before(time.Now().Add(-i.cfg.MaxAge)) {
			logx.Debugf("skipping message because date is not in range (msg.UID=%d)", msg.UID)
			continue
		}

		for {
			p, err = mr.NextPart()
			if errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				logx.Warnf("failed to read message part (msg.UID=%d): %v\n", msg.UID, err)
				break
			}

			switch p.Header.(type) {
			case *mail.InlineHeader:
				if b, err = io.ReadAll(p.Body); err != nil {
					logx.Warnf("failed to read message body (msg.UID=%d): %v\n", msg.UID, err)
					break
				}
				message.Contents = append(message.Contents, string(b))
			}
		}

		messages = append(messages, message)
	}

	return messages, nil
}

func (i *Imap) MoveMessage(uid imap.UID, mailbox string) error {
	uidSet := imap.UIDSet{}
	uidSet.AddNum(uid)
	if _, err := i.client.Move(uidSet, mailbox).Wait(); err != nil {
		return err
	}
	return nil
}
