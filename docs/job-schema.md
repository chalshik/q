# Job Schema

## Job

|Field|Type|Description|
|---|---|---|
|`id`|string (UUID)|Unique job identifier|
|`user_id`|string|ID of the producer that submitted the job|
|`prompt`|string|The prompt to be executed|
|`status`|enum|Current job state (see below)|
|`created_at`|timestamp|When the job was submitted|
|`updated_at`|timestamp|When the job status last changed|

---

## Job Status

|Status|Description|
|---|---|
|`queued`|Job received and waiting in queue|
|`running`|Job assigned to a worker and executing|
|`done`|Job completed successfully|
|`failed`|Job failed due to worker error or timeout|

---

## State Machine

```
queued → running → done
                 ↘ failed
```

### Transitions

|From|To|Trigger|Owner|
|---|---|---|---|
|—|`queued`|Job received via REST|Networking|
|`queued`|`running`|Worker pulls job|Scheduler|
|`running`|`done`|Worker reports success|Worker|
|`running`|`failed`|Worker reports error|Worker|
|`running`|`failed`|Heartbeat timeout|Scheduler|