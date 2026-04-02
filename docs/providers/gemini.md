# Gemini

Uses Gemini to analyze the message.

Configuration options:

| Field     | Type    | Required | Description                                      | Example            |
|-----------|---------|----------|--------------------------------------------------|--------------------|
| `apikey`  | string  | yes      | OpenAI API key                                   | `some-api-key`     |
| `model`   | string  | yes      | OpenAI model used for classification             | `gemini-2.5-flash` |
| `maxsize` | integer | yes      | Maximum email size sent to the model (bytes)     | `100000`           |
| `prompt`  | string  | no       | The prompt which is sent to the model (optional) | _see above_        |

Example:

```yaml
providers:
  prov1:
    type: gemini
    config:
      apikey: some-api-key
      model: gemini-2.5-flash
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
