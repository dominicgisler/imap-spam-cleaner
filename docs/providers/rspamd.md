# Rspamd

Uses a Rspamd server for classification.

The tool uses the `/checkv2` endpoint to classify the message.
The returned `score` and `required_score` will be used to calculate the internal `score` (percentage).
If you want to rely on the `required_score` of `rspamd`, set the `minscore` option in your inbox to `100`.

Configuration options:

| Field     | Type     | Required | Description                                  | Example                  |
|-----------|----------|----------|----------------------------------------------|--------------------------|
| `url`     | string   | yes      | Rspamd API URL                               | `http://127.0.0.1:11333` |
| `timeout` | duration | no       | Timeout for the request                      | `10s`                    |

Example:

```yaml
providers:
  prov1:
    type: rspamd
    config:
      url: http://127.0.0.1:11333
      timeout: 10s
```
