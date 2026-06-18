# API Contracts

## REST API (Producer → Dispatcher)

### Submit a job

```
POST /jobs
```

**Request body**

```json
{
  "user_id": "string",
  "prompt": "string"
}
```

**Response**

```json
{
  "id": "string",
  "status": "queued",
  "created_at": "timestamp"
}
```

---

### Get job status

```
GET /jobs/:id
```

**Response**

```json
{
  "id": "string",
  "user_id": "string",
  "status": "queued | running | done | failed",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

---

## gRPC API (Worker → Dispatcher)

### PullJob

Worker requests the next available job.

**Request**

```protobuf
message PullJobRequest {
  string worker_id = 1;
}
```

**Response**

```protobuf
message PullJobResponse {
  string job_id = 1;
  string prompt = 2;
}
```

---

### ReportResult

Worker reports job completion or failure.

**Request**

```protobuf
message ReportResultRequest {
  string worker_id = 1;
  string job_id   = 2;
  string status   = 3; // "done" or "failed"
  string result   = 4; // output or error message
}
```

**Response**

```protobuf
message ReportResultResponse {
  bool accepted = 1;
}
```

---

### Heartbeat

Worker signals it is still alive.

**Request**

```protobuf
message HeartbeatRequest {
  string worker_id = 1;
}
```

**Response**

```protobuf
message HeartbeatResponse {
  bool acknowledged = 1;
}
```