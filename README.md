# IMAP spam cleaner

[![Docker Hub](https://img.shields.io/badge/Docker%20Hub-grey?logo=docker)](https://hub.docker.com/r/dominicgisler/imap-spam-cleaner)
[![License](https://img.shields.io/github/license/dominicgisler/imap-spam-cleaner)](https://github.com/dominicgisler/imap-spam-cleaner/blob/master/LICENSE)
[![Issues](https://img.shields.io/github/issues/dominicgisler/imap-spam-cleaner)](https://github.com/dominicgisler/imap-spam-cleaner/issues)
[![Last commit](https://img.shields.io/github/last-commit/dominicgisler/imap-spam-cleaner/master)](https://github.com/dominicgisler/imap-spam-cleaner/commits/master)

A tool to clean up spam in your imap inbox.

## How does it work

This application loads mails from configured imap inboxes and checks their contents using the defined provider.
Depending on a spam score, the message can be moved to the spam folder, keeping your inbox clean.

## How to use

### Using image from docker hub (recommended)

- Create `config.yml` matching your inboxes (example below)
- Create `docker-compose.yml` if using `docker compose` (example below)
- Start the container with: `docker compose up -d`
- Or with: `docker run -d --name imap-spam-cleaner -v ./config.yml:/app/config.yml dominicgisler/imap-spam-cleaner`

### From source with local Go installation

- Clone this repository
- Install Go version 1.23.4+
- Load dependencies (`go get ./...`)
- Create `config.yml` matching your inboxes (example below)
- Run the application (`go run .`)

### From source with docker

- Clone this repository
- Install docker
- Build the docker image: `docker build -f Dockerfile -t dominicgisler/imap-spam-cleaner .`
- Create `config.yml` matching your inboxes (example below)
- Create `docker-compose.yml` if using `docker compose` (example below)
- Start the container with: `docker compose up -d`
- Or with: `docker run -d --name imap-spam-cleaner -v ./config.yml:/app/config.yml dominicgisler/imap-spam-cleaner`

### Sample docker-compose.yml

```yaml
services:
  imap-spam-cleaner:
    image: dominicgisler/imap-spam-cleaner:latest
    container_name: imap-spam-cleaner
    hostname: imap-spam-cleaner
    restart: always
    volumes:
      - ./config.yml:/app/config.yml:ro
```

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
