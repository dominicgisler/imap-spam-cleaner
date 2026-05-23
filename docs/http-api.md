# HTTP API

IMAP Spam Cleaner exposes HTTP endpoints for runtime information when the HTTP server is started.

The data is stored in `store.db`, make sure to map the file accordingly on your system.

## Run summary

```http
GET /runs/summary
```

Returns one summary entry per inbox. Each entry contains the first recorded run, the last recorded run, the number of runs, and the summed run counters.

### Query parameters

| Name    | Required | Description                         |
|---------|----------|-------------------------------------|
| `inbox` | no       | Limits the summary to one inbox.    |

### Example request unfiltered

```console
curl "http://localhost:8080/runs/summary"
```

```json
[
  {
    "inbox": "test@example.com",
    "run_count": 10,
    "first_run_at": "2026-05-23 08:00:00+00:00",
    "last_run_at": "2026-05-23 11:55:00+00:00",
    "message_count": 100,
    "skipped_count": 10,
    "failed_count": 2,
    "moved_count": 88
  },
  {
    "inbox": "info@example.com",
    "run_count": 20,
    "first_run_at": "2026-05-23 08:00:00+00:00",
    "last_run_at": "2026-05-23 11:55:00+00:00",
    "message_count": 200,
    "skipped_count": 20,
    "failed_count": 4,
    "moved_count": 176
  }
]
```

### Example request for a specific inbox

```console
curl "http://localhost:8080/runs/summary?inbox=info@example.com"
```

```json
{
  "inbox": "info@example.com",
  "run_count": 12,
  "first_run_at": "2026-05-23 08:00:00+00:00",
  "last_run_at": "2026-05-23 11:55:00+00:00",
  "message_count": 184,
  "skipped_count": 9,
  "failed_count": 2,
  "moved_count": 37
}
```

When `inbox` is omitted, an array of summaries will be returned.
If an inbox is provided, only one summary element will be returned.

### Fields

| Field           | Description                                      |
|-----------------|--------------------------------------------------|
| `inbox`         | Inbox identifier used by the run.                |
| `run_count`     | Number of stored runs for the inbox.             |
| `first_run_at`  | Start time of the first stored run.              |
| `last_run_at`   | Start time of the most recent stored run.        |
| `message_count` | Total number of messages loaded across runs.     |
| `skipped_count` | Total number of skipped messages across runs.    |
| `failed_count`  | Total number of failed message analyses.         |
| `moved_count`   | Total number of messages moved to the spam folder. |

Only the `GET` method is supported. Other methods return `405 Method Not Allowed`.
