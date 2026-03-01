# Cron Jobs (Scheduled Tasks)

Schedule recurring tasks with the Goose cron module.

## Overview

The Cron module enables time-based job scheduling:

- Run tasks at specific times
- Periodic data cleanup
- Report generation
- Health checks
- External API synchronization

## Quick Start

```go
import "github.com/awesome-goose/goose/modules/cron"

// Configure cron module
cronModule := cron.NewModule()

// Include in application
stop, err := goose.Start(goose.API(platform, module, []types.Module{
    cronModule,
}))
```

## Defining Scheduled Tasks

### Basic Cron Job

```go
package jobs

import "github.com/awesome-goose/goose/modules/cron"

type CleanupJob struct {
    db *gorm.DB `inject:""`
}

func (j *CleanupJob) Handle() error {
    // Delete old records
    return j.db.Where("created_at < ?", time.Now().AddDate(0, -6, 0)).Delete(&Log{}).Error
}

// Run every day at midnight
func (j *CleanupJob) Schedule() string {
    return "0 0 * * *"
}
```

### Cron Expression Format

```
┌───────────── minute (0 - 59)
│ ┌───────────── hour (0 - 23)
│ │ ┌───────────── day of month (1 - 31)
│ │ │ ┌───────────── month (1 - 12)
│ │ │ │ ┌───────────── day of week (0 - 6) (Sunday = 0)
│ │ │ │ │
* * * * *
```

### Common Schedules

```go
// Every minute
func (j *Job) Schedule() string { return "* * * * *" }

// Every 5 minutes
func (j *Job) Schedule() string { return "*/5 * * * *" }

// Every hour
func (j *Job) Schedule() string { return "0 * * * *" }

// Every day at midnight
func (j *Job) Schedule() string { return "0 0 * * *" }

// Every day at 9 AM
func (j *Job) Schedule() string { return "0 9 * * *" }

// Every Monday at 9 AM
func (j *Job) Schedule() string { return "0 9 * * 1" }

// First day of every month at midnight
func (j *Job) Schedule() string { return "0 0 1 * *" }

// Every weekday at 6 PM
func (j *Job) Schedule() string { return "0 18 * * 1-5" }
```

## Registering Cron Jobs

### In Module

```go
func (m *AppModule) Declarations() []any {
    return []any{
        // Register cron jobs
        &CleanupJob{},
        &DailyReportJob{},
        &HealthCheckJob{},
    }
}
```

### Direct Registration

```go
cronModule := cron.NewModule(
    cron.WithJobs(
        cron.Job{
            Name:     "cleanup",
            Schedule: "0 0 * * *",
            Handler: func() error {
                return cleanupOldData()
            },
        },
        cron.Job{
            Name:     "health-check",
            Schedule: "*/5 * * * *",
            Handler: func() error {
                return checkHealth()
            },
        },
    ),
)
```

## Job Configuration

### With Timezone

```go
type ReportJob struct{}

func (j *ReportJob) Handle() error {
    return generateDailyReport()
}

func (j *ReportJob) Schedule() string {
    return "0 9 * * *" // 9 AM
}

func (j *ReportJob) Timezone() string {
    return "America/New_York"
}
```

### Single Instance (No Overlap)

```go
type LongRunningJob struct{}

func (j *LongRunningJob) Handle() error {
    // This job may take longer than schedule interval
    return processLargeDataset()
}

func (j *LongRunningJob) Schedule() string {
    return "*/5 * * * *"
}

// Don't start new instance if previous is still running
func (j *LongRunningJob) WithoutOverlapping() bool {
    return true
}
```

### With Timeout

```go
func (j *Job) Timeout() time.Duration {
    return 5 * time.Minute
}
```

## Practical Examples

### Database Cleanup

```go
type DatabaseCleanupJob struct {
    db *gorm.DB `inject:""`
}

func (j *DatabaseCleanupJob) Handle() error {
    // Delete expired sessions
    if err := j.db.Where("expires_at < ?", time.Now()).Delete(&Session{}).Error; err != nil {
        return err
    }

    // Delete old logs
    if err := j.db.Where("created_at < ?", time.Now().AddDate(0, -1, 0)).Delete(&Log{}).Error; err != nil {
        return err
    }

    // Delete unverified users older than 24 hours
    if err := j.db.Where("verified = ? AND created_at < ?", false, time.Now().Add(-24*time.Hour)).Delete(&User{}).Error; err != nil {
        return err
    }

    return nil
}

func (j *DatabaseCleanupJob) Schedule() string {
    return "0 2 * * *" // 2 AM daily
}
```

### Daily Report Generation

```go
type DailyReportJob struct {
    reportService *ReportService `inject:""`
    emailService  *EmailService  `inject:""`
}

func (j *DailyReportJob) Handle() error {
    // Generate report for yesterday
    yesterday := time.Now().AddDate(0, 0, -1)
    report, err := j.reportService.GenerateDailyReport(yesterday)
    if err != nil {
        return err
    }

    // Send to admins
    admins := j.reportService.GetAdminEmails()
    for _, email := range admins {
        j.emailService.SendReport(email, report)
    }

    return nil
}

func (j *DailyReportJob) Schedule() string {
    return "0 6 * * *" // 6 AM daily
}
```

### Health Check

```go
type HealthCheckJob struct {
    services      []string
    alertService  *AlertService `inject:""`
}

func (j *HealthCheckJob) Handle() error {
    for _, service := range j.services {
        if err := j.checkService(service); err != nil {
            j.alertService.SendAlert("Service unhealthy: " + service)
        }
    }
    return nil
}

func (j *HealthCheckJob) checkService(service string) error {
    resp, err := http.Get(service + "/health")
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return fmt.Errorf("unhealthy: status %d", resp.StatusCode)
    }
    return nil
}

func (j *HealthCheckJob) Schedule() string {
    return "*/5 * * * *" // Every 5 minutes
}
```

### Cache Warming

```go
type CacheWarmingJob struct {
    cache          *cache.Client   `inject:""`
    productService *ProductService `inject:""`
}

func (j *CacheWarmingJob) Handle() error {
    // Warm popular products cache
    products, _ := j.productService.GetPopularProducts(100)
    for _, product := range products {
        j.cache.Set("product:"+product.ID, product, 2*time.Hour)
    }

    // Warm categories cache
    categories, _ := j.productService.GetCategories()
    j.cache.Set("categories:all", categories, 2*time.Hour)

    return nil
}

func (j *CacheWarmingJob) Schedule() string {
    return "0 */2 * * *" // Every 2 hours
}
```

### External API Sync

```go
type SyncExternalDataJob struct {
    apiClient    *ExternalAPI   `inject:""`
    dataService  *DataService   `inject:""`
}

func (j *SyncExternalDataJob) Handle() error {
    // Fetch external data
    data, err := j.apiClient.FetchLatestData()
    if err != nil {
        return fmt.Errorf("failed to fetch external data: %w", err)
    }

    // Update local database
    return j.dataService.UpdateFromExternal(data)
}

func (j *SyncExternalDataJob) Schedule() string {
    return "0 */4 * * *" // Every 4 hours
}
```

## Error Handling

### With Error Callback

```go
type ImportJob struct{}

func (j *ImportJob) Handle() error {
    return runImport()
}

func (j *ImportJob) OnError(err error) {
    log.Printf("Import job failed: %v", err)
    sendAlertEmail("Import job failed", err.Error())
}

func (j *ImportJob) Schedule() string {
    return "0 1 * * *"
}
```

### With Retry

```go
type RetryableJob struct {
    attempts int
}

func (j *RetryableJob) Handle() error {
    j.attempts++
    err := doWork()

    if err != nil && j.attempts < 3 {
        // Will be retried
        return err
    }

    return nil
}

func (j *RetryableJob) MaxRetries() int {
    return 3
}
```

## Monitoring Cron Jobs

### Job History

```go
type CronJobLog struct {
    ID        string    `json:"id"`
    JobName   string    `json:"job_name"`
    StartedAt time.Time `json:"started_at"`
    EndedAt   time.Time `json:"ended_at"`
    Status    string    `json:"status"` // success, failed
    Error     string    `json:"error,omitempty"`
}

func (j *Job) Handle() error {
    log := &CronJobLog{
        ID:        uuid.New().String(),
        JobName:   "MyJob",
        StartedAt: time.Now(),
    }

    err := j.doWork()

    log.EndedAt = time.Now()
    if err != nil {
        log.Status = "failed"
        log.Error = err.Error()
    } else {
        log.Status = "success"
    }

    j.saveLog(log)
    return err
}
```

### Status Endpoint

```go
func (c *AdminController) CronStatus(ctx types.Context) any {
    jobs := c.cronService.GetJobs()

    status := make([]map[string]interface{}, len(jobs))
    for i, job := range jobs {
        status[i] = map[string]interface{}{
            "name":      job.Name,
            "schedule":  job.Schedule,
            "last_run":  job.LastRun,
            "next_run":  job.NextRun,
            "status":    job.Status,
        }
    }

    return status
}
```

## Best Practices

1. **Make jobs idempotent** - Safe to run multiple times
2. **Use appropriate schedules** - Don't run too frequently
3. **Handle errors gracefully** - Log and alert on failures
4. **Set timeouts** - Prevent hung jobs
5. **Avoid overlapping** for long-running jobs
6. **Monitor execution** - Track success/failure rates
7. **Use distributed locks** for multi-instance deployments

## Next Steps

- [Queues](queues.md) - Background job processing
- [Events](events.md) - Event-driven architecture
- [Logging](../building-blocks/logging.md) - Job logging
