package imap

import (
	"fmt"
	"github.com/dominicgisler/imap-spam-cleaner/config"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
)

func CheckInbox(i config.Inbox) error {

	var err error
	var client *imapclient.Client
	var mbox *imap.SelectData
	var msgs []*imapclient.FetchMessageBuffer

	if i.TLS {
		client, err = imapclient.DialTLS(fmt.Sprintf("%s:%d", i.Host, i.Port), nil)
	} else {
		client, err = imapclient.DialInsecure(fmt.Sprintf("%s:%d", i.Host, i.Port), nil)
	}
	if err != nil {
		return fmt.Errorf("failed to dial IMAP server: %v", err)
	}
	defer func() { _ = client.Close() }()

	if err = client.Login(i.Username, i.Password).Wait(); err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	mbox, err = client.Select(i.Inbox, nil).Wait()
	if err != nil {
		return fmt.Errorf("failed to select INBOX: %v", err)
	}

	if mbox.NumMessages > 0 {
		seqSet := imap.SeqSet{}
		seqSet.AddRange(1, mbox.NumMessages)
		fetchOptions := &imap.FetchOptions{Envelope: true}
		msgs, err = client.Fetch(seqSet, fetchOptions).Collect()
		if err != nil {
			return fmt.Errorf("failed to fetch first message in INBOX: %v", err)
		}
		for _, msg := range msgs {
			fmt.Println(msg.Envelope.Subject)
		}
	}

	if err = client.Logout().Wait(); err != nil {
		return fmt.Errorf("failed to logout: %v", err)
	}

	return nil
}
