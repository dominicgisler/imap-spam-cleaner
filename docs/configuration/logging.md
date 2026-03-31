# Logging

| Field   | Type   | Required | Default | Description         | Options                                                      |
|---------|--------|----------|---------|---------------------|--------------------------------------------------------------|
| `level` | string | no       | info    | Log verbosity level | `panic`, `fatal`, `error`, `warn`, `info`, `debug`, `trace ` |

Defining a level will enable all levels above and including the used one.
Example: setting it to `warn` will enable `panic`, `fatal`, `error` and `warn` - using `trace` will enable all possible logging messages.

At the moment only `error`, `warn`, `info` and `debug` are in use.

Example:

```yaml
logging:
  level: debug
```
