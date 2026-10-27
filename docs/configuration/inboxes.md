# Inboxes

| Field       | Type     | Required | Description                                                                                         | Example           |
|-------------|----------|----------|-----------------------------------------------------------------------------------------------------|-------------------|
| `schedule`  | string   | yes      | Cron schedule defining when spam analysis runs, see [Runs](#runs)                                   | `*/5 * * * *`     |
| `host`      | string   | yes      | IMAP server hostname                                                                                | `mail.domain.tld` |
| `port`      | integer  | yes      | IMAP port, usually `993` with TLS                                                                   | `993`             |
| `tls`       | boolean  | no       | Connect with TLS, `false` connects unencrypted (STARTTLS is not supported)                          | `true`            |
| `username`  | string   | yes      | IMAP username                                                                                       | `user@domain.tld` |
| `password`  | string   | yes      | IMAP password                                                                                       | `mypass`          |
| `provider`  | string   | yes      | Provider used for spam detection                                                                    | `prov1`           |
| `inbox`     | string   | yes      | Folder to scan                                                                                      | `INBOX`           |
| `spam`      | string   | yes      | Folder where spam messages are moved                                                                | `INBOX.Spam`      |
| `markread`  | boolean  | no       | Mark detected messages as read                                                                      | `true`            |
| `minscore`  | integer  | yes      | Minimum spam score required to classify as spam (0–100), see [scores](../providers/index.md#scores) | `75`              |
| `minage`    | duration | no       | Minimum age of message before scanning                                                              | `0h`              |
| `maxage`    | duration | no       | Maximum age of message considered, empty = no limit, see [Runs](#runs)                              | `15m`             |
| `whitelist` | string   | no       | [Whitelist](whitelists.md) to use                                                                   | `whitelist1`      |

Durations are written like `90s`, `15m` or `24h`, days are not supported.

With `tls: false`, the password and all messages are sent in plain text.
Only use it for a server on the same host or in a trusted network.

## Runs

Each run analyzes every message in `inbox` that is older than `minage` and younger than `maxage`, based on the time the server received it.
Results are not remembered between runs: a message that isn't moved is analyzed again on every run for as long as it stays in that range,
and with a paid AI provider every analysis is billed. Without `maxage`, every message in the folder is analyzed on every run.

Keep `maxage` a little longer than the schedule interval. For example, `*/5 * * * *` with `maxage: 15m` analyzes each message about three times.
Messages that arrive while the application is stopped for longer than `maxage` are never analyzed.

Runs of the same inbox never overlap. If a run is still in progress when the next one is due, that run is skipped.

Schedules use the standard 5-field cron syntax and are evaluated in UTC in the Docker image.

Example:

```yaml
inboxes:
  - schedule: "*/5 * * * *"
    host: mail.domain.tld
    port: 993
    tls: true
    username: user@domain.tld
    password: mypass
    provider: prov1
    inbox: INBOX
    spam: INBOX.Spam
    markread: true
    minscore: 75
    minage: 0h
    maxage: 15m
    whitelist: whitelist1
```
