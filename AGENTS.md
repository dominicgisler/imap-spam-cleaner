# AGENTS.md

This file provides guidance to AI coding agents working with code in this repository.

## Overview

Go service that fetches mail from configured IMAP inboxes on a cron schedule, scores each message 0–100 with a pluggable spam provider (OpenAI, Ollama, Gemini, SpamAssassin, rspamd), and moves messages scoring `>= minscore` to the spam folder. User docs live in `docs/` (MkDocs, deployed to GitHub Pages on push to master).

## Commands

```bash
go run .                     # needs ./config.yml (copy config.example.yml); creates ./store.db
go run . -now                # process every inbox once, ignoring schedules, then exit (no HTTP server)
go run . -analyze-only       # score only: never mark/move messages, never record runs
make lint                    # golangci-lint; .golangci.yml sets issues.fix, so it rewrites files
make fix                     # go fix ./... (modernizers)
make docs                    # preview docs via Docker on :8000
make TAG=x.y.z               # multi-arch Docker build (version injected via -ldflags main.version)
docker compose -f docker-compose.dev.yml up   # local spamd (:783) and rspamd (:11333)
```

There are no automated tests yet. Any new ones should use the standard `testing` package (`go test ./...`; single test: `go test ./provider -run TestName`).

CI: golangci-lint v2.15.0 with `only-new-issues` (linters: errorlint, gocritic, gocyclo, gosec, usestdlibvars; gofmt/goimports). golangci-lint refuses to run if it was built with an older Go than the `go` directive in `go.mod`, so bump the pin in `.github/workflows/golangci-lint.yml` (and the local install) together with the Go version. Docs CI runs `mkdocs build --strict`, so a broken link or nav entry fails the deploy.

## Architecture

`main.go` loads config → `database.Init("store.db")` → validates every provider config up front (fail fast) → either `inbox.RunAllInboxes` (`-now`) or starts the HTTP API on `:8080` in a goroutine and blocks in `inbox.Schedule` (gocron, one cron job per inbox, stops on SIGINT/SIGTERM). The paths `config.yml` and `store.db` and the port are hardcoded and resolved against the working directory (`/app` in Docker).

`inbox.processInbox` is the core loop: open an IMAP connection, load messages, skip senders on the whitelist, `provider.Analyze` each message one at a time, optionally mark `\Seen`, move it to the spam folder, then write a `database.Run` row. Analysis failures are counted and skipped, never fatal. There is **no per-message state**: every run re-scores every message still in the folder within the `minage`/`maxage` window, so the schedule and `maxage` directly drive LLM API cost.

- **Config (`app/`)**: YAML is checked with go-playground/validator tags, then cross-checked: an inbox must name an existing provider and whitelist. `required` rejects zero values, so don't put it on ints or bools where `0`/`false` is valid (it's why `minscore: 0` is currently rejected). Whitelist entries decode straight into `regexp.Regexp`.
- **IMAP (`imap/`)**: fetches the full body with `Peek` (doesn't set `\Seen`). IMAP `SINCE`/`BEFORE` search only works by day, so `LoadMessages` re-filters precisely on `INTERNALDATE` on the client side. `Message.Contents` holds only inline parts (attachments are dropped). `Raw` and `RawHeader` hold the original bytes. Whitelist regexes match the raw `From` header value (`Name <addr@host>`), which is why the examples look like `^.* <addr>$`.
- **Providers (`provider/`)**: the `Provider` interface has `Name`, `ValidateConfig`, `Init` and `Analyze`. Provider config is a flat `map[string]string`, so each provider parses its own keys. `ValidateConfig` stores the parsed values on the receiver and `Init` calls it before building a client. `provider.New` returns a fresh instance per inbox run. Scores are normalised differently per provider, so `minscore` means something different for each:
  - **AI** (`AIBase`, embedded by openai/ollama/gemini): the prompt is a `text/template` (vars: `From`, `To`, `DeliveredTo`, `Cc`, `Bcc`, `Subject`, `Content`, `Raw`, `RawHeader`). The reply goes through `strconv.ParseInt` without trimming, so any extra text or whitespace counts as a failed analysis. `maxsize` drops whole content parts once their combined size passes the limit; `Raw`/`RawHeader` are never cut.
  - **SpamAssassin**: speaks the spamd `REPORT` protocol (SPAMC/1.5) over raw TCP. The score is mapped linearly from -2…10 to 0…100 and clamped. `maxsize` truncates the raw bytes.
  - **rspamd**: POSTs the raw message to `/checkv2` and returns `round(score / required_score * 100)`, unclamped, so 100 equals rspamd's own threshold. There is no `maxsize`.
- **Database (`database/`)**: SQLite through pure-Go `modernc.org/sqlite` (the Docker build uses `CGO_ENABLED=0`) and sqlx. The SQL is embedded from `database/queries/` and loaded with `query("run/create")`. Goose migrations live in `queries/migrations`. `migrateDB` migrates up **or down** to exactly `migVersion`, so bump that number with every new migration.
- **HTTP (`handler/`)**: only `GET /runs/summary?inbox=&maxage=<seconds>` (see `docs/http-api.md`).
- **Logging (`logx/`)**: a logrus wrapper. Use it instead of `log`/`fmt`.

The module targets Go 1.27, so newer syntax such as `new(expr)` (in `database/run.go`) is intentional and valid.

### Adding a provider

Touch all of these:
- a new file in `provider/` and an entry in the `provider.New` list
- the `oneof=` tag on `app.Provider.Type`
- `config.example.yml` and the full config example in `docs/getting-started.md`
- `docs/providers/<name>.md`, the list in `docs/providers/index.md`, and the `mkdocs.yml` nav

## Contribution conventions (CONTRIBUTING.md)

- Open an issue before large changes. Keep each PR to one feature or fix. Lint must pass.
- Commit messages are short; `fix:`/`feat:`/`docs:` prefixes are optional.
- Update `docs/` whenever behaviour or config changes.
- AI tools may be used, but contributors remain responsible for the quality and correctness of everything they submit and must understand and review it.
- Releases are cut manually as GitHub Releases. A `vX.Y.Z` tag publishes Docker images tagged `vX.Y.Z`, `vX.Y`, `vX` and `latest`.
