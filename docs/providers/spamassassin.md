# SpamAssassin

Uses a SpamAssassin server for classification.

The message is checked with spamd's `REPORT` command. SpamAssassin's score is mapped from -2…10 to 0…100 (lower and higher values are clamped),
so the internal `score` is `(score + 2) / 12 * 100`. SpamAssassin's default threshold of 5.0 equals a `minscore` of 58, `minscore: 75` moves messages scoring 7.0 or more.

Configuration options:

| Field     | Type     | Required | Description                                                                | Example     |
|-----------|----------|----------|----------------------------------------------------------------------------|-------------|
| `host`    | string   | no       | SpamAssassin host, default `127.0.0.1`                                     | `127.0.0.1` |
| `port`    | integer  | no       | SpamAssassin port, default `783`                                           | `783`       |
| `maxsize` | integer  | no       | Maximum message size sent to SpamAssassin (bytes), longer messages are cut | `300000`    |
| `timeout` | duration | no       | Timeout for connecting and for the response, default `5s`                  | `10s`       |

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
