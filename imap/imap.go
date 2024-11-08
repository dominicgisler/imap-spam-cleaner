package imap

import (
	"bytes"
	"fmt"
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	_ "github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	"io"
)

func LoadMessages(i config.Inbox) ([]Message, error) {

	var err error
	var client *imapclient.Client
	var mbox *imap.SelectData
	var msgs []*imapclient.FetchMessageBuffer
	var mr *mail.Reader
	var p *mail.Part
	var messages []Message

	if i.TLS {
		client, err = imapclient.DialTLS(fmt.Sprintf("%s:%d", i.Host, i.Port), nil)
	} else {
		client, err = imapclient.DialInsecure(fmt.Sprintf("%s:%d", i.Host, i.Port), nil)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to dial IMAP server: %v", err)
	}
	defer func() { _ = client.Close() }()

	if err = client.Login(i.Username, i.Password).Wait(); err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	mbox, err = client.Select(i.Inbox, nil).Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %v", err)
	}

	if mbox.NumMessages > 0 {
		seqSet := imap.SeqSet{}
		seqSet.AddRange(1, mbox.NumMessages)
		fetchOptions := &imap.FetchOptions{
			Envelope:    true,
			BodySection: []*imap.FetchItemBodySection{{Peek: true}},
		}
		msgs, err = client.Fetch(seqSet, fetchOptions).Collect()
		if err != nil {
			return nil, fmt.Errorf("failed to fetch messages: %v", err)
		}

		for _, msg := range msgs {
			var b []byte
			for _, buf := range msg.BodySection {
				b = buf
				break
			}

			mr, err = mail.CreateReader(bytes.NewReader(b))
			if err != nil {
				return nil, fmt.Errorf("failed to create mail reader: %v", err)
			}

			message := Message{
				DeliveredTo: mr.Header.Get("Delivered-To"),
				From:        mr.Header.Get("From"),
				To:          mr.Header.Get("To"),
				Cc:          mr.Header.Get("Cc"),
				Bcc:         mr.Header.Get("Bcc"),
				Subject:     msg.Envelope.Subject,
				Contents:    []string{},
			}

			for {
				p, err = mr.NextPart()
				if err == io.EOF {
					break
				} else if err != nil {
					return nil, fmt.Errorf("failed to read next part: %v", err)
				}

				switch p.Header.(type) {
				case *mail.InlineHeader:
					if b, err = io.ReadAll(p.Body); err != nil {
						return nil, fmt.Errorf("failed to create mail body: %v", err)
					}
					message.Contents = append(message.Contents, string(b))
				}
			}

			messages = append(messages, message)
		}
	}

	if err = client.Logout().Wait(); err != nil {
		return nil, fmt.Errorf("failed to logout: %v", err)
	}

	return messages, nil
}
