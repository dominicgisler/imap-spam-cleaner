package provider

import (
	"fmt"

	"github.com/dominicgisler/imap-spam-cleaner/imap"
	"github.com/dominicgisler/imap-spam-cleaner/logx"
)

type AIBase struct {
	model   string
	maxsize int
	prompt  string
}

func (p *AIBase) buildPrompt(msg imap.Message) string {

	cont := ""
	contLen := 0
	for _, cnt := range msg.Contents {
		contLen += len(cnt)
		if contLen > p.maxsize {
			logx.Debugf("skipping bytes for message #%d (%s)", msg.UID, msg.Subject)
			break
		}
		cont += cnt + "\n"
	}

	return fmt.Sprintf(
		"Analyze the following email for its spam potential.\n"+
			"Return a spam score between 0 and 100. Only answer with the number itself.\n\n"+
			"From: %s\nTo: %s\nDelivered-To: %s\nCc: %s\nBcc: %s\nSubject: %s\n\n"+
			"Content:\n%s",
		msg.From,
		msg.To,
		msg.DeliveredTo,
		msg.Cc,
		msg.Bcc,
		msg.Subject,
		cont,
	)
}
