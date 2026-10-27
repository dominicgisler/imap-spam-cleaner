# Providers

- [OpenAI](openai.md)
- [Ollama](ollama.md)
- [Gemini](gemini.md)
- [TypeSafe (Jev)](typesafe.md)
- [SpamAssassin](spamassassin.md)
- [Rspamd](rspamd.md)

## Scores

Scores are on a scale from 0 to 100 (rspamd's can be negative or exceed 100). A message is moved if its score is at least the inbox's `minscore`.
The score means something different for each provider:

| Provider               | Score                                                                               |
|------------------------|-------------------------------------------------------------------------------------|
| OpenAI, Ollama, Gemini | The number the model answers, as asked for in the prompt                            |
| TypeSafe               | Probability of spam in percent                                                      |
| SpamAssassin           | SpamAssassin's score mapped from -2…10 to 0…100, its default threshold of 5.0 is 58 |
| Rspamd                 | Percentage of rspamd's `required_score`, `100` equals rspamd's own threshold        |

Messages are analyzed again on every run, see [how runs work](../configuration/inboxes.md#runs).

## AI

Each request to an AI provider is cancelled after `timeout` (default `2m`).
The message then counts as failed and is analyzed again on the next run.

The following prompt is used by the AI providers if no custom prompt is specified.

```
Analyze the following email for its spam potential.
Return a spam score between 0 and 100. Only answer with the number itself.

{{.RawHeader}}

Content:
{{.Content}}
```

The following placeholders can be used in the prompt:

- `{{.From}}`: message header `From`
- `{{.To}}`: message header `To`
- `{{.DeliveredTo}}`: message header `Delivered-To`
- `{{.Cc}}`: message header `Cc`
- `{{.Bcc}}`: message header `Bcc`
- `{{.Subject}}`: message subject
- `{{.Content}}`: message content (cut after `maxsize`)
- `{{.Raw}}`: full raw message (not affected by `maxsize`)
- `{{.RawHeader}}`: all message headers (not affected by `maxsize`)
