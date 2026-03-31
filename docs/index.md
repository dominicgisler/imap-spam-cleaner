# IMAP Spam Cleaner

![logo](assets/icon_128.png)

IMAP Spam Cleaner is a tool that automatically cleans spam from IMAP inboxes.

The application scans messages using a configured provider and moves spam messages to a spam folder based on their spam score.

## Features

- Works with any IMAP mailbox
- Spam detection via configurable providers
- Runs as Docker container
- Supports multiple inboxes
- YAML configuration

## Example

```console
$ docker run -v ./config.yml:/app/config.yml dominicgisler/imap-spam-cleaner:latest
INFO   [2026-02-28T16:53:41Z] Starting imap-spam-cleaner v0.5.3
DEBUG  [2026-02-28T16:53:41Z] Loaded config
INFO   [2026-02-28T16:53:41Z] Scheduling inbox info@example.com (*/5 * * * *)
INFO   [2026-02-28T16:55:00Z] Handling info@example.com
DEBUG  [2026-02-28T16:55:00Z] Available mailboxes:
DEBUG  [2026-02-28T16:55:00Z]   - INBOX
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Drafts
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Sent
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Trash
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Spam
DEBUG  [2026-02-28T16:55:00Z]   - INBOX.Spam.Cleaner
DEBUG  [2026-02-28T16:55:00Z] Found 34 messages in inbox
DEBUG  [2026-02-28T16:55:00Z] Found 5 UIDs in timerange
INFO   [2026-02-28T16:55:00Z] Loaded 5 messages
DEBUG  [2026-02-28T16:55:06Z] Spam score of message #478 (Herzlichen Glückwunsch! Ihr Decathlon-Geschenk wartet auf Sie. 🎁): 90/100
DEBUG  [2026-02-28T16:55:12Z] Spam score of message #479 (Leider ist bei der Verarbeitung Ihrer Zahlung ein Problem aufgetreten.): 90/100
DEBUG  [2026-02-28T16:55:18Z] Spam score of message #480 (Das neue Geheimnis gegen Bauchfett!): 92/100
DEBUG  [2026-02-28T16:55:26Z] Spam score of message #481 (Schnell: 1 Million / Lady Million): 80/100
DEBUG  [2026-02-28T16:55:32Z] Spam score of message #483 (Vermögen x4 zu Fest): 85/100
INFO   [2026-02-28T16:55:32Z] Moved 4 messages
```
