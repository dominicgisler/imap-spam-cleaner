package imap

type Message struct {
	DeliveredTo string
	From        string
	To          string
	Cc          string
	Bcc         string
	Subject     string
	Contents    []string
}
