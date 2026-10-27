# TypeSafe (Jev)

Uses [TypeSafe](https://typesafe.ai)'s Jev model to classify the message.

Jev is a classifier, not a text generator: the message is sent together with a yes/no question
(`instructions`) and Jev returns the probability that the answer is yes.
That probability is used as the internal `score` (percentage), so `minscore: 75` moves every message
that is at least 75% likely to be spam.
Only input tokens are billed.

Configuration options:

| Field          | Type    | Required | Description                                                                 | Example                                   |
|----------------|---------|----------|-----------------------------------------------------------------------------|-------------------------------------------|
| `apikey`       | string  | yes      | TypeSafe API key                                                            | `some-api-key`                            |
| `model`        | string  | yes      | Jev model used for classification                                           | `jev-1.14.0`                              |
| `maxsize`      | integer | no       | Maximum message size sent to the model (bytes), only affects `{{.Content}}` | `50000`                                   |
| `instructions` | string  | no       | The yes/no question asked about the message                                 | `Is this email spam, phishing or a scam?` |
| `prompt`       | string  | no       | The message sent to the model (optional), the question is sent separately   | _see below_                               |
| `url`          | string  | no       | API URL, e.g. for a proxy                                                   | `https://api.typesafe.ai`                 |

`jev-latest` always points to the newest release. Pin a version such as `jev-1.14.0` to keep scores stable.
The available models are listed in the [TypeSafe docs](https://docs.typesafe.ai/models).

A request may use at most 32k tokens per question, so set `maxsize` to keep long messages from being rejected.
Rejected messages are logged and skipped.

If no `prompt` is set, the following one is used. It supports the same placeholders as the
[other AI providers](index.md#ai), but must not contain a question, use `instructions` for that.

```
{{.RawHeader}}

Content:
{{.Content}}
```

Example:

```yaml
providers:
  prov1:
    type: typesafe
    config:
      apikey: some-api-key
      model: jev-1.14.0
      maxsize: 50000
      instructions: Is this email spam, phishing or a scam?
      prompt: |
        From: {{.From}}
        To: {{.To}}
        Subject: {{.Subject}}

        Content:
        {{.Content}}
```
