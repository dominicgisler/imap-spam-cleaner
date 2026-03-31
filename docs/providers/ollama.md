# Ollama

Uses a local LLM to analyze the message.

Configuration options:

| Field     | Type    | Required | Description                                      | Example                  |
|-----------|---------|----------|--------------------------------------------------|--------------------------|
| `url`     | string  | yes      | Ollama server URL                                | `http://127.0.0.1:11434` |
| `model`   | string  | yes      | Ollama model name used for classification        | `gpt-oss:20b`            |
| `maxsize` | integer | yes      | Maximum email size sent to the model (bytes)     | `100000`                 |
| `prompt`  | string  | no       | The prompt which is sent to the model (optional) | _see above_              |

Example:

```yaml
providers:
  prov1:
    type: ollama
    config:
      url: http://127.0.0.1:11434
      model: gpt-oss:20b
      maxsize: 100000
      prompt: |
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
```
