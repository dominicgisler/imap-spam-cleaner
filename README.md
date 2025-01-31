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

- Install Go version 1.23.4+
- Clone this repository
- Load dependencies (`go get ./...`)
- Create `config.yml` matching your inboxes
- Run the application (`go run .`)

### With docker

- Build the docker image: `docker build -f Dockerfile -t dominicgisler/imap-spam-cleaner .`
- Start the container with: `docker compose up -d`
- Or with: `docker run -d --name imap-spam-cleaner -v ./config.yml:/app/config.yml dominicgisler/imap-spam-cleaner`

### Configuration

Use this configuration as an example for your own setup. Save the file as `config.yml` on your disk (where the application will run) or mount the correct path into the docker container.

```yaml
logging:
  level: debug                # logging level (panic, fatal, error, warn, info, debug, trace)

providers:                    # providers to be used for inboxes
  prov1:                      # provider name
    type: openai              # provider type
    config:                   # provider specific configuration
      apikey: some-api-key    # apikey
      model: gpt-4o-mini      # gpt model to use
      maxsize: 100000         # message size limit for prompt (bytes)

inboxes:                      # inboxes to be checked
  - schedule: "* * * * *"     # schedule in cron format (when to execute spam analysis)
    host: mail.domain.tld     # imap host
    port: 143                 # imap port
    tls: false                # imap tls
    username: user@domain.tld # imap user
    password: mypass          # imap password
    provider: prov1           # provider used for spam analysis
    inbox: INBOX              # inbox folder
    spam: INBOX.Spam          # spam folder
    minscore: 75              # min score to detect spam (0-100)
    minage: 0h                # min age of message
    maxage: 24h               # max age of message
```
