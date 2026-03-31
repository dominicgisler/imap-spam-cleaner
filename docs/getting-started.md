# Getting started

- Create `config.yml` matching your inboxes (example below)
- Create `docker-compose.yml` if using docker compose (example below)
- Start the container with: `docker compose up -d`
- Or with: `docker run -d --name imap-spam-cleaner -v ./config.yml:/app/config.yml dominicgisler/imap-spam-cleaner`

The container will run in the background and execute analysis steps according to the defined schedule. If needed check logs either with `docker compose logs -f` or `docker logs -f imap-spam-cleaner`.

## Configuration

Use this configuration as an example, it contains different providers with different configuration options. Consult specific wiki pages for detailed information about the options.

```yaml
logging:
  level: debug                    # logging level (panic, fatal, error, warn, info, debug, trace)

providers:                        # providers to be used for inboxes
  prov1:                          # provider name
    type: openai                  # provider type
    config:                       # provider specific configuration
      apikey: some-api-key        # openai apikey
      model: gpt-4o-mini          # openai model to use
      maxsize: 100000             # message size limit for prompt (bytes)
      prompt: |                   # prompt to be sent to the model
        Analyze the following email for its spam potential.
        Return a spam score between 0 and 100. Only answer with the number itself.

        From: {{.From}}
        To: {{.To}}
        Delivered-To: {{.DeliveredTo}}
        Cc: {{.Cc}}
        Bcc: {{.Bcc}}
        Subject: {{.Subject}}

        Content:
        {{.Content}}
  prov2:                          # provider name
    type: ollama                  # provider type
    config:                       # provider specific configuration
      url: http://127.0.0.1:11434 # ollama url
      model: gpt-oss:20b          # ollama model to use
      maxsize: 100000             # message size limit for prompt (bytes)
  prov3:                          # provider name
    type: spamassassin            # provider type
    config:                       # provider specific configuration
      host: 127.0.0.1             # spamassassin host
      port: 783                   # spamassassin port
      maxsize: 300000             # message size limit

inboxes:                          # inboxes to be checked
  - schedule: "* * * * *"         # schedule in cron format (when to execute spam analysis)
    host: mail.domain.tld         # imap host
    port: 143                     # imap port
    tls: false                    # imap tls
    username: user@domain.tld     # imap user
    password: mypass              # imap password
    provider: prov1               # provider used for spam analysis
    inbox: INBOX                  # inbox folder
    spam: INBOX.Spam              # spam folder
    minscore: 75                  # min score to detect spam (0-100)
    minage: 0h                    # min age of message
    maxage: 24h                   # max age of message
```

## Docker compose

This compose-file shows a minimal setup, the `config.yml` can be anywhere on your system, but needs to be mapped to `/app/config.yml`.

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