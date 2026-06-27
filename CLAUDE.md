# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

**Start the database:**
```bash
docker compose up -d
```

**Run the dispatcher:**
```bash
go run ./cmd/dispatcher
```

**Run the worker:**
```bash
go run ./cmd/worker
```

**Run all tests:**
```bash
go test ./...
```

**Run a single test:**
```bash
go test ./internal/networking/http/... -run TestCreateJobHandler
```

**Regenerate protobuf code** (requires `protoc` and `protoc-gen-go`):
```bash
protoc --go_out=. --go-grpc_out=. proto/dispatcher.proto
```

## Architecture

The system is a prompt dispatcher that sits in front of LLM inferencworkers. It guarantees exactly-once execution under worker failure by tracking job state in Postgres and auto-requeueing timed-out jobs.e 

```
Producer → REST → Dispatcher → gRPC → Workers
                      |
                  Postgres
```

**Two binaries:**
- `cmd/dispatcher` — the central server: accepts jobs from producers, manages job lifecycle, serves workers over gRPC.
- `cmd/worker` — a pull-based worker: loops forever pulling jobs from the dispatcher via gRPC, executes them, and submits results.

### Internal layers (`internal/`)

- **`networking/http`** — REST handler for producers (`POST /jobs`). Creates jobs via `JobService`.
- **`networking/grpc`** — gRPC server for workers. `GetJob` pops a job ID from the queue and fetches it; `SubmitResult` updates job status.
- **`services`** — `JobService` is the single business logic entry point: assigns UUIDs, sets initial status, writes to Postgres, and pushes to the in-memory queue.
- **`scheduler`** — Runs a background loop that pops from the queue and drives job dispatch. Intended to also handle heartbeat timeouts and worker failure detection.
- **`queue`** — A bounded in-memory channel queue. Decoupled behind its own type so it can be replaced (e.g. priority queue) without touching the scheduler.
- **`persistence`** — All Postgres access. `NewPostgresConnection` reads env vars, runs embedded migrations automatically on startup, and returns a `*sqlx.DB`. `JobRepository` owns all SQL queries.
- **`models`** — Shared `Job` struct used across all layers.

### Database

Migrations live in `db/migrations/` as `.up.sql`/`.down.sql` files and are embedded via `db/migrations.go`. They run automatically when the dispatcher starts — no manual migration step needed.

**Connection env vars** (defaults work with `docker-compose.yaml`):
| Var | Default |
|---|---|
| `DB_HOST` | `localhost` |
| `DB_PORT` | `5432` |
| `DB_USER` | `postgres` |
| `DB_PASSWORD` | `secret` |
| `DB_NAME` | `dispatcher` |
| `DB_SSLMODE` | `disable` |

### Job lifecycle

```
queued → running → done
                 ↘ failed
```

Status transitions are owned by specific layers: `queued` is set on REST ingestion, `running` when a worker pulls, `done`/`failed` when the worker reports back.

### Proto

`proto/dispatcher.proto` defines two RPCs: `GetJob` and `SubmitResult`. The generated Go files (`dispatcher.pb.go`, `dispatcher_grpc.pb.go`) are checked in — only regenerate if the `.proto` changes.
