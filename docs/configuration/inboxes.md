# Inboxes

| Field       | Type     | Required | Description                                             | Example           |
|-------------|----------|----------|---------------------------------------------------------|-------------------|
| `schedule`  | string   | yes      | Cron schedule defining when spam analysis runs          | `*/5 * * * *`     |
| `host`      | string   | yes      | IMAP server hostname                                    | `mail.domain.tld` |
| `port`      | integer  | yes      | IMAP port                                               | `143`             |
| `tls`       | boolean  | no       | Enable TLS                                              | `false`           |
| `username`  | string   | yes      | IMAP username                                           | `user@domain.tld` |
| `password`  | string   | yes      | IMAP password                                           | `mypass`          |
| `provider`  | string   | yes      | Provider used for spam detection                        | `prov1`           |
| `inbox`     | string   | yes      | Folder to scan                                          | `INBOX`           |
| `spam`      | string   | yes      | Folder where spam messages are moved                    | `INBOX.Spam`      |
| `minscore`  | integer  | yes      | Minimum spam score required to classify as spam (0–100) | `75`              |
| `minage`    | duration | no       | Minimum age of message before scanning                  | `0h`              |
| `maxage`    | duration | no       | Maximum age of message considered                       | `24h`             |
| `whitelist` | string   | no       | Whitelist to use                                        | `whitelist1`      |

```yaml
inboxes:
  - schedule: "* * * * *"
    host: mail.domain.tld
    port: 143
    tls: false
    username: user@domain.tld
    password: mypass
    provider: prov1
    inbox: INBOX
    spam: INBOX.Spam
    minscore: 75
    minage: 0h
    maxage: 24h
    whitelist: whitelist1
```
