# Providers

- [OpenAI](openai.md)
- [Ollama](ollama.md)
- [Gemini](gemini.md)
- [SpamAssassin](spamassassin.md)

## AI

The following prompt is used by the AI providers if no custom prompt is specified.

```
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
