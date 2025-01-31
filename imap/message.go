package imap

import (
	"github.com/emersion/go-imap/v2"
	"time"
)

type Message struct {
	UID         imap.UID
	DeliveredTo string
	From        string
	To          string
	Cc          string
	Bcc         string
	Subject     string
	Contents    []string
	Date        time.Time
}
