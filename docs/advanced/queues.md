# Background Queues

Process background jobs asynchronously with the Goose `queues` module.

## Overview

The Queues module is a SQL-backed job queue with an auto-scaling worker
pool, optional retries, exponential backoff, and stale-job recovery. Use it
for: sending emails, processing uploads, generating reports, external API
calls, and other resource-intensive operations.

The queue persists in the SQL connection registered by the SQL module — it
is **not** Redis-backed.

## Quick Start

```go
import (
    "context"

    "github.com/awesome-goose/goose/modules/queues"
    "github.com/awesome-goose/goose/modules/sql"
    "github.com/awesome-goose/goose/types"
)

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        sql.Root(&sql.Config{Dialect: "sqlite", Name: "app.db", Sync: true}),
        queues.Root(queues.DefaultConfig(), m.queueHandlers()...),
    }
}

func (m *AppModule) queueHandlers() []*queues.JobHandler {
    return []*queues.JobHandler{
        queues.NewSimpleHandler("emails", "send-welcome", func(j *queues.QueueJob) (any, error) {
            data, err := queues.GetJobData[WelcomeEmail](j)
            if err != nil {
                return nil, err
            }
            return nil, sendWelcomeEmail(data.UserID)
        }).WithMinWorkers(1).WithMaxWorkers(4),
    }
}
```

`queues.Root(cfg, handlers...)` is shorthand for
`queues.NewModule(cfg, handlers, true)`; `queues.Child(handlers...)` reuses
the root queue from a child module.

## Configuration

```go
type Config struct {
    Queue                  string        // Default queue name for new jobs
    DefaultRetryLimit      int           // Retries per failed job (default 3)
    DefaultRetryDelay      int           // ms between retries
    PollInterval           time.Duration // Worker poll cadence
    CleanupInterval        time.Duration // Sweep frequency (0 = never)
    RetentionPeriod        time.Duration // Keep completed/failed jobs (default 7d)
    StaleJobTimeout        time.Duration // Recover jobs stuck in_progress (default 5m)
    EnableStaleJobRecovery bool          // Auto-recover stale jobs
}

cfg := queues.DefaultConfig()
queues.WithQueue("default")(cfg)
queues.WithPollInterval(5 * time.Second)(cfg)
queues.WithDefaultRetryLimit(3)(cfg)
```

There is **no `Workers` field** and no `WithDriver`/`WithWorkers` option —
worker scaling is per-handler via `WithMinWorkers`/`WithMaxWorkers`.

## Handlers

A `*queues.JobHandler` describes how to run a particular `queue:job`. It
includes the function plus worker-pool settings:

```go
type JobHandler struct {
    Queue                 string         // queue name (e.g. "emails")
    Job                   string         // job type (e.g. "send-welcome")
    Handler               JobHandlerFn   // func(ctx, *QueueJob) (any, error)
    MinWorkers, MaxWorkers int
    PollIntervalMs        int
    IdleThreshold         int            // ticks of inactivity before scale-down
    TimeoutMs             int            // per-execution timeout (default 60_000)
    UseExponentialBackoff bool
}
```

### Constructors

```go
// Plain handler
h := queues.NewHandler("emails", "send-welcome",
    func(ctx context.Context, j *queues.QueueJob) (any, error) {
        return nil, deliver(j)
    })

// Simple handler (no ctx)
h := queues.NewSimpleHandler("emails", "send-welcome",
    func(j *queues.QueueJob) (any, error) {
        return nil, deliver(j)
    })

// Fluent worker config
h.WithMinWorkers(1).WithMaxWorkers(8).
   WithPollInterval(500).
   WithIdleThreshold(10).
   WithTimeout(30_000).
   WithExponentialBackoff()
```

### Typed handlers

`queues.NewTypedHandler[T]` decodes the persisted JSON payload for you:

```go
type SendEmail struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
}

th := queues.NewTypedHandler[SendEmail]("emails", "send",
    func(ctx context.Context, data SendEmail, job *queues.QueueJob) (any, error) {
        return nil, sendEmail(data.To, data.Subject)
    })

handlers := []*queues.JobHandler{th.ToJobHandler()}
```

## Pushing Jobs

Inject `*queues.Queue` and call one of the push methods:

```go
type OrderController struct {
    queue *queues.Queue `inject:""`
}

func (c *OrderController) Create(ctx types.Context) any {
    var dto CreateOrderDTO
    _ = ctx.Bind(&dto)

    // Push immediately to the default queue
    job, _ := c.queue.Push("emails", "send-welcome",
        SendEmail{To: dto.Email, Subject: "Welcome"}, nil)

    return map[string]string{"job_id": job.Id}
}
```

### Push with delay or at a specific time

```go
// Push 5 minutes from now
job, _ := c.queue.Delay("emails", "reminder", payload, 5*time.Minute, nil)

// Push at an absolute time
job, _ = c.queue.At("reports", "monthly", payload, time.Date(2026, 1, 1, 9, 0, 0, 0, time.UTC), nil)

// Fluent later-style API
job, _ = c.queue.Later("emails", "reminder", payload, nil).InMinutes(5)
```

### Push a batch

```go
batch := []queues.BatchJob{
    {Name: "send", Data: SendEmail{To: "a@example.com"}},
    {Name: "send", Data: SendEmail{To: "b@example.com"}},
}
jobs, _ := c.queue.PushBatch("emails", batch)
```

## Per-Job Configuration

Pass a `*queues.JobConfig` to `Push`/`Delay`/`At` to override defaults for a
single job:

```go
priority := 10
c.queue.Push("emails", "send-welcome", payload, &queues.JobConfig{
    Priority:   priority,
    RetryLimit: 5,
    RetryDelay: 5_000, // ms
    TimeoutMs:  30_000,
    Singleton:  true,  // skip if another job with same name is pending
})
```

## Error Handling and Retries

When a handler returns an error, the runner retries up to
`RetryLimit + 1` attempts. Between retries it sleeps `RetryDelay` ms — or
doubles the delay (capped) if the handler enabled `WithExponentialBackoff`.
Jobs that exhaust their retries are recorded with `status = "failed"`.

You can re-queue a failed job manually:

```go
err := c.queue.Retry(jobID)
```

## Monitoring

```go
m := c.queue.GetMetrics()
// m.TotalProcessed, m.TotalSucceeded, m.TotalFailed, m.TotalRetried, m.ProcessingTimes

job, _   := c.queue.GetJob(jobID)
logs, _  := c.queue.GetJobLogs(jobID)
jobs, _  := c.queue.ListJobs("emails", queues.JobStatusNew, 100)

_ = c.queue.PauseQueue("emails")
_ = c.queue.ResumeQueue("emails")
_ = c.queue.Cancel(jobID)
```

## Practical Examples

### Welcome email

```go
type WelcomeEmail struct {
    UserID string `json:"user_id"`
}

h := queues.NewTypedHandler[WelcomeEmail]("emails", "send-welcome",
    func(ctx context.Context, data WelcomeEmail, job *queues.QueueJob) (any, error) {
        user, err := userService.GetByID(data.UserID)
        if err != nil {
            return nil, err
        }
        return nil, sendWelcomeEmail(user.Email, user.Name)
    }).ToJobHandler().WithMinWorkers(2).WithMaxWorkers(8)
```

### File processing

```go
h := queues.NewSimpleHandler("uploads", "process",
    func(j *queues.QueueJob) (any, error) {
        return nil, processUpload(j) // resize, convert, etc.
    }).
    WithTimeout(30 * 60 * 1000). // 30 min
    WithMaxWorkers(2)
```

### Report generation

```go
h := queues.NewSimpleHandler("reports", "generate",
    func(j *queues.QueueJob) (any, error) {
        return nil, runReport(j)
    }).
    WithMaxWorkers(1) // serialize heavy jobs
```

## Best Practices

1. **Keep jobs small and idempotent** — they may run more than once.
2. **Set a sensible `TimeoutMs`** on long-running handlers.
3. **Separate queues** by latency requirement (e.g. `emails` vs `reports`).
4. **Use `Singleton`** for jobs that must not run concurrently.
5. **Enable `EnableStaleJobRecovery`** in multi-instance deployments.
6. **Watch `GetMetrics()`** for failure spikes; expose an admin endpoint.

## What this module does NOT provide

- No Redis backend. Jobs persist via the SQL module.
- No `Dispatch().Delay()/At()/OnQueue()` chained DSL. Use `Push`, `Delay`,
  or `At` directly.
- No `Chain`/`Batch.Then` higher-order composition primitives — chain via
  pushing the next job from inside the previous handler.
- The type is `*queues.Queue`, not `*queues.Client`.

## Next Steps

- [Cron Jobs](cron.md) - Scheduled tasks
- [Caching](caching.md) - Performance optimization
- [Events](events.md) - Event-driven architecture
