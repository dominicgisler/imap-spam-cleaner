# SpamAssassin

Uses a SpamAssassin server for classification.

Configuration options:

| Field     | Type     | Required | Description                                  | Example     |
|-----------|----------|----------|----------------------------------------------|-------------|
| `host`    | string   | yes      | SpamAssassin host                            | `127.0.0.1` |
| `port`    | integer  | yes      | SpamAssassin port                            | `783`       |
| `maxsize` | integer  | no       | Maximum email size sent to the model (bytes) | `300000`    |
| `timeout` | duration | no       | Timeout for the request                      | `10s`       |

Example:

```yaml
providers:
  prov1:
    type: spamassassin
    config:
      host: 127.0.0.1
      port: 783
      maxsize: 1000000
      timeout: 10s
```
