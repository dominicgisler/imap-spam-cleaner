package imap

import (
	"time"

	"github.com/emersion/go-imap/v2"
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
	TextBody    string
	HtmlBody    string
	Date        time.Time
	Raw         []byte
	Headers     string
}
