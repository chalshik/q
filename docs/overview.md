# Dispatcher Overview

## What is it?

A prompt dispatcher that sits in front of LLM inference workers. It accepts prompts from producers, queues them, and distributes them to available workers via a pull-based model.

---

## High Level Architecture

```
Producer → REST → Dispatcher → gRPC → Workers
                      |
                  Database
```

**Producer** — submits jobs (prompt + metadata) via REST API.

**Dispatcher** — central node that manages job lifecycle. Contains Networking, Scheduler, Queue, and Persistence layers.

**Workers** — inference servers that pull jobs when free, execute the prompt, and report results back.

**Database** — persistent store for job state.

---

## Dispatcher Internals

### Networking Layer

Responsible for two interfaces:

- **Inbound** — exposes REST API for producers to submit jobs
- **Outbound** — exposes gRPC methods for workers to pull jobs, send heartbeats, and report results

### Scheduler Layer

Manages worker registration and tracks worker availability. Applies FIFO scheduling to select the next job when a worker pulls. Detects worker failures via heartbeat timeout and marks affected jobs as failed.

### Queue Layer

Owned by the Scheduler. Holds pending jobs in order of arrival. Decoupled behind an interface to allow future replacement (e.g. priority queue, delay queue) without changing Scheduler logic.

### Persistence Layer

Manages all database operations. Single access point for job state transitions. Called by:

- **Networking** — on job received (status: `queued`)
- **Worker** — on job completion (status: `done` or `failed`)

---

## Data Flow

```
1. Producer submits job via REST
2. Networking validates and writes job to Queue (status: queued)
3. Worker sends gRPC pull request to Networking → Scheduler
4. Scheduler picks next job from Queue, marks it running, returns to Worker
5. Worker executes prompt
6. Worker writes result to Persistence (status: done/failed)
7. Worker signals Scheduler it is free
```