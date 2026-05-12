# Cron Jobs (Scheduled Tasks)

Schedule recurring tasks with the Goose `cron` module.

## Overview

The Cron module runs handlers on a cron schedule. Each handler is registered
under a `group:name` pair, persisted in the SQL backend, and picked up by a
runner that ticks at a configurable interval.

Use it for: periodic data cleanup, report generation, health checks, and
external API synchronization.

## Quick Start

```go
import (
    "context"

    "github.com/awesome-goose/goose/modules/cron"
    "github.com/awesome-goose/goose/modules/sql"
    "github.com/awesome-goose/goose/types"
)

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        sql.Root(&sql.Config{Dialect: "sqlite", Name: "app.db", Sync: true}),
        cron.NewModule(&cron.Config{
            TickInterval:           time.Minute,
            Timezone:               "Etc/UTC",
            EnableStaleJobRecovery: true,
        }, m.cronHandlers(), false),
    }
}

func (m *AppModule) cronHandlers() []*cron.CronHandler {
    return []*cron.CronHandler{
        cron.NewSimpleHandler("system", "cleanup", "0 2 * * *", func(j *cron.CronJob) (any, error) {
            return nil, runDailyCleanup()
        }),
    }
}
```

`cron.NewModule(config, handlers, isRoot)` takes the handler slice
explicitly. The runner is started by the framework — your module just
provides the handlers.

## Handlers

A `*cron.CronHandler` is the unit of scheduling:

```go
type CronHandler struct {
    Group                 string                // logical group (e.g. "billing")
    Name                  string                // unique within the group
    Pattern               string                // cron expression
    Handler               CronHandlerFn         // func(ctx, *CronJob) (any, error)
    Config                *CronConfig           // optional priority/retry settings
    TimeoutMs             int                   // execution timeout (default 60_000)
    UseExponentialBackoff bool                  // retry backoff strategy
}
```

### Constructors

```go
// With ctx-aware handler:
h := cron.NewHandler("system", "cleanup", "0 2 * * *",
    func(ctx context.Context, j *cron.CronJob) (any, error) {
        return runCleanup(ctx)
    })

// Without ctx (convenience):
h := cron.NewSimpleHandler("system", "ping", "*/5 * * * *",
    func(j *cron.CronJob) (any, error) {
        return nil, ping()
    })
```

### Fluent configuration

```go
h := cron.NewHandler("billing", "daily-report", "0 6 * * *", reportFn).
    WithPriority(10).
    WithRetryLimit(3).
    WithRetryDelay(5_000).             // ms (capped at 15_000)
    WithTimeout(10 * 60 * 1000).       // 10 minute timeout
    WithExponentialBackoff()
```

### Typed handlers

Use `cron.NewTypedHandler[T]` when each job carries a JSON config payload:

```go
type ReportArgs struct {
    Region string `json:"region"`
}

th := cron.NewTypedHandler[ReportArgs]("billing", "regional-report", "0 6 * * *",
    func(ctx context.Context, args ReportArgs, job *cron.CronJob) (any, error) {
        return generateReport(args.Region)
    })

handlers := []*cron.CronHandler{th.ToCronHandler()}
```

## Cron Expression Format

```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday = 0)
│ │ │ │ │
* * * * *
```

Common patterns:

```go
"* * * * *"       // every minute
"*/5 * * * *"     // every 5 minutes
"0 * * * *"       // every hour
"0 0 * * *"       // every day at midnight
"0 9 * * 1"       // Mondays at 9 AM
"0 0 1 * *"       // first day of the month
"0 18 * * 1-5"    // weekdays at 6 PM
```

## Config

```go
type Config struct {
    TickInterval           time.Duration // how often the runner wakes up (default 60s)
    Timezone               string        // IANA name (default "Etc/UTC")
    DefaultRetryLimit      int           // 3
    DefaultRetryDelay      int           // milliseconds
    StaleJobTimeout        time.Duration // 5 minute default
    EnableStaleJobRecovery bool          // recover jobs stuck in_progress
    CleanupInterval        time.Duration // log retention sweep (0 = never)
    LogRetentionPeriod     time.Duration // 7 days default
}
```

Functional setters:

```go
cfg := cron.DefaultConfig().Apply(
    cron.WithTimezone("America/New_York"),
    cron.WithTickInterval(30*time.Second),
    cron.WithStaleJobRecovery(true),
)
```

## Practical Examples

### Database cleanup

```go
type Cleanup struct {
    db *sql.Db `inject:""`
}

func (c *Cleanup) Handler(ctx context.Context, j *cron.CronJob) (any, error) {
    cutoff := time.Now().AddDate(0, -1, 0)
    return nil, c.db.Where("created_at < ?", cutoff).Delete(&Log{}).Error
}

func (m *AppModule) cronHandlers() []*cron.CronHandler {
    cleanup := &Cleanup{}
    return []*cron.CronHandler{
        cron.NewHandler("system", "log-cleanup", "0 2 * * *", cleanup.Handler),
    }
}
```

### Cache warming

```go
type Warmer struct {
    cache   *cache.Cache    `inject:""`
    service *ProductService `inject:""`
}

func (w *Warmer) Handler(ctx context.Context, j *cron.CronJob) (any, error) {
    products, _ := w.service.GetPopularProducts(100)
    for _, p := range products {
        _ = w.cache.Set("product:"+p.ID, p, 2*time.Hour)
    }
    return len(products), nil
}
```

### Health check

```go
type Health struct {
    alerts *AlertService `inject:""`
    urls   []string
}

func (h *Health) Handler(ctx context.Context, j *cron.CronJob) (any, error) {
    for _, u := range h.urls {
        resp, err := http.Get(u + "/health")
        if err != nil || resp.StatusCode != 200 {
            h.alerts.SendAlert("unhealthy: " + u)
        }
        if resp != nil {
            resp.Body.Close()
        }
    }
    return nil, nil
}
```

## Error Handling and Retries

When the handler returns an error, the runner retries up to
`RetryLimit + 1` attempts (the initial run plus `RetryLimit` retries).
Between retries it sleeps `RetryDelay` ms — or, if
`UseExponentialBackoff: true`, doubles the delay each attempt up to
`DefaultMaxBackoffDelay` (60s).

```go
h := cron.NewHandler("imports", "nightly", "0 1 * * *", runImport).
    WithRetryLimit(3).
    WithRetryDelay(5_000).
    WithExponentialBackoff()
```

## Monitoring

Each execution writes a `CronLog` row keyed by `job_id`. Query the runner's
in-memory metrics via `(*cron.Cron).GetMetrics()`:

```go
m := cronService.GetMetrics()
// m.TotalRuns, m.TotalSucceeded, m.TotalFailed, m.TotalRetried, m.LastRunAt
```

Or list jobs and recent logs:

```go
jobs, _ := cronService.ListJobs("")          // all groups
logs, _ := cronService.GetJobLogs(jobID, 20) // last 20 entries
```

## Best Practices

1. **Make handlers idempotent** — they may run more than once on retry.
2. **Pick the smallest schedule that meets the need** — every minute adds DB
   churn.
3. **Set a `TimeoutMs`** for long-running jobs.
4. **Enable `EnableStaleJobRecovery`** in multi-instance deployments.
5. **Use `Priority`** to keep critical jobs ahead of bulk work.

## Next Steps

- [Queues](queues.md) - Ad-hoc background job processing
- [Events](events.md) - Event-driven architecture
- [Logging](../building-blocks/logging.md) - Job logging
