# IMAP spam cleaner

[![License](https://img.shields.io/github/license/dominicgisler/imap-spam-cleaner)](https://github.com/dominicgisler/imap-spam-cleaner/blob/master/LICENSE)
[![Issues](https://img.shields.io/github/issues/dominicgisler/imap-spam-cleaner)](https://github.com/dominicgisler/imap-spam-cleaner/issues)
[![Last commit](https://img.shields.io/github/last-commit/dominicgisler/imap-spam-cleaner/master)](https://github.com/dominicgisler/imap-spam-cleaner/commits/master)

A tool to clean up spam in your imap inbox.

**Work In Progress**

## How does it work

This application loads mails from configured imap inboxes and checks their contents using the defined provider.
Depending on a spam score, the message can be moved to the spam folder, keeping your inbox clean.

## How to use

### From source

- Install Go version 1.23.2+
- Clone this repository
- Load dependencies (`go get ./...`)
- Create `config.yml` matching your inboxes
- Run the application (`go run .`)

### Configuration

Use this configuration as an example for your own setup.

```yaml
providers:                    # providers to be used for inboxes
  prov1:                      # provider name
    type: openai              # provider type
    credentials:              # provider specific credentials
      apikey: some-api-key

inboxes:                      # inboxes to be checked
  - host: mail.domain.tld     # imap host
    port: 143                 # imap port
    tls: false                # imap tls
    username: user@domain.tld # imap user
    password: mypass          # imap password
    provider: prov1           # provider used for spam analysis
    inbox: INBOX              # inbox folder
    spam: INBOX.Spam          # spam folder
    minscore: 75              # min score to detect spam (0-100)
```
