# Providers

| Field    | Type   | Required | Description                     | Options                            |
|----------|--------|----------|---------------------------------|------------------------------------|
| `type`   | string | yes      | Provider implementation         | `openai`, `ollama`, `spamassassin` |
| `config` | object | yes      | Provider-specific configuration |                                    |

A list of providers, which can be reused by inboxes.
Each provider is identified by a name (`prov1` in the following example) and will be referenced by that name.
You can find a detailed description for the options on the individual provider pages.

Example:

```yaml
providers:
  prov1:
    type: openai
    config:
      apikey: some-api-key
      model: gpt-4o-mini
      maxsize: 100000
```
