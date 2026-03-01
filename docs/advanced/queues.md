# Background Queues

Process background jobs asynchronously with the Goose queues module.

## Overview

The Queues module enables background job processing for:

- Sending emails
- Processing uploads
- Generating reports
- External API calls
- Resource-intensive operations

## Quick Start

```go
import "github.com/awesome-goose/goose/modules/queues"

// Configure queues module
queuesModule := queues.NewModule(
    queues.WithDriver("redis"),
    queues.WithHost("localhost"),
    queues.WithPort(6379),
    queues.WithQueues("default", "emails", "exports"),
)

// Include in application
stop, err := goose.Start(goose.API(platform, module, []types.Module{
    queuesModule,
}))
```

## Defining Jobs

### Basic Job

```go
package jobs

import (
    "github.com/awesome-goose/goose/modules/queues"
)

type SendEmailJob struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
    Body    string `json:"body"`
}

func (j *SendEmailJob) Handle(ctx queues.JobContext) error {
    // Send the email
    return sendEmail(j.To, j.Subject, j.Body)
}

func (j *SendEmailJob) Queue() string {
    return "emails"
}
```

### Job with Dependencies

```go
type ProcessOrderJob struct {
    OrderID string `json:"order_id"`
}

func (j *ProcessOrderJob) Handle(ctx queues.JobContext) error {
    // Get services from context
    orderService := ctx.Get("orderService").(*OrderService)

    // Process the order
    return orderService.Process(j.OrderID)
}

func (j *ProcessOrderJob) Queue() string {
    return "orders"
}

func (j *ProcessOrderJob) Timeout() time.Duration {
    return 5 * time.Minute
}
```

## Dispatching Jobs

### Inject Queue Client

```go
type OrderController struct {
    queue        *queues.Client   `inject:""`
    orderService *OrderService    `inject:""`
}
```

### Dispatch Immediately

```go
func (c *OrderController) Create(ctx types.Context) any {
    var dto CreateOrderDTO
    ctx.Bind(&dto)

    // Create order
    order, err := c.orderService.Create(dto)
    if err != nil {
        return ctx.Status(500).JSON(map[string]string{"error": err.Error()})
    }

    // Dispatch background job
    c.queue.Dispatch(&ProcessOrderJob{OrderID: order.ID})

    return ctx.Status(201).JSON(order)
}
```

### Dispatch with Delay

```go
// Process 5 minutes from now
c.queue.Dispatch(&SendReminderJob{
    UserID: user.ID,
}).Delay(5 * time.Minute)

// Process at specific time
c.queue.Dispatch(&SendReportJob{
    ReportID: report.ID,
}).At(time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC))
```

### Dispatch to Specific Queue

```go
c.queue.Dispatch(&SendEmailJob{
    To:      user.Email,
    Subject: "Welcome!",
    Body:    "Welcome to our platform.",
}).OnQueue("emails")
```

## Job Configuration

### Retries

```go
type ProcessPaymentJob struct {
    PaymentID string `json:"payment_id"`
}

func (j *ProcessPaymentJob) Handle(ctx queues.JobContext) error {
    return processPayment(j.PaymentID)
}

// Retry 3 times on failure
func (j *ProcessPaymentJob) MaxRetries() int {
    return 3
}

// Wait between retries
func (j *ProcessPaymentJob) RetryDelay() time.Duration {
    return 30 * time.Second
}
```

### Timeout

```go
func (j *ExportDataJob) Timeout() time.Duration {
    return 10 * time.Minute
}
```

### Priority

```go
func (j *UrgentNotificationJob) Priority() int {
    return 10 // Higher priority
}

func (j *RegularNotificationJob) Priority() int {
    return 1 // Default priority
}
```

## Error Handling

### Handle Failures

```go
type SendEmailJob struct {
    To      string `json:"to"`
    Subject string `json:"subject"`
}

func (j *SendEmailJob) Handle(ctx queues.JobContext) error {
    err := sendEmail(j.To, j.Subject)
    if err != nil {
        return fmt.Errorf("failed to send email: %w", err)
    }
    return nil
}

// Called when all retries exhausted
func (j *SendEmailJob) OnFailure(ctx queues.JobContext, err error) {
    log.Printf("Email to %s failed after retries: %v", j.To, err)
    // Store failed job for manual review
    storeFailedJob(j, err)
}
```

### Conditional Retry

```go
func (j *APICallJob) ShouldRetry(err error) bool {
    // Don't retry on validation errors
    if errors.Is(err, ErrValidation) {
        return false
    }
    // Retry on network errors
    return true
}
```

## Job Chaining

Execute jobs in sequence:

```go
func (c *Controller) ProcessOrder(ctx types.Context) any {
    orderID := ctx.Param("id")

    // Chain jobs
    c.queue.Chain(
        &ValidateOrderJob{OrderID: orderID},
        &ProcessPaymentJob{OrderID: orderID},
        &SendConfirmationJob{OrderID: orderID},
        &UpdateInventoryJob{OrderID: orderID},
    ).Dispatch()

    return map[string]string{"status": "processing"}
}
```

## Job Batching

Process multiple jobs together:

```go
func (c *Controller) SendBulkEmails(ctx types.Context) any {
    var dto BulkEmailDTO
    ctx.Bind(&dto)

    // Create batch
    jobs := make([]queues.Job, len(dto.Recipients))
    for i, recipient := range dto.Recipients {
        jobs[i] = &SendEmailJob{
            To:      recipient,
            Subject: dto.Subject,
            Body:    dto.Body,
        }
    }

    // Dispatch batch
    batch := c.queue.Batch(jobs...).OnQueue("emails")

    // Optional: callback when all complete
    batch.Then(&BatchCompleteJob{BatchID: batch.ID})

    batch.Dispatch()

    return map[string]string{"batch_id": batch.ID}
}
```

## Worker Configuration

### Register Job Handlers

```go
// In your module
func (m *AppModule) Declarations() []any {
    return []any{
        // Register jobs
        &SendEmailJob{},
        &ProcessOrderJob{},
        &ExportDataJob{},
    }
}
```

### Configure Workers

```go
queuesModule := queues.NewModule(
    queues.WithDriver("redis"),
    queues.WithQueues("default", "emails", "exports"),
    queues.WithWorkers(map[string]int{
        "default": 2,  // 2 workers for default queue
        "emails":  4,  // 4 workers for emails
        "exports": 1,  // 1 worker for exports (sequential)
    }),
)
```

## Monitoring Jobs

### Job Status

```go
func (c *Controller) GetJobStatus(ctx types.Context) any {
    jobID := ctx.Param("id")

    status, err := c.queue.GetStatus(jobID)
    if err != nil {
        return ctx.Status(404).JSON(map[string]string{"error": "Job not found"})
    }

    return map[string]interface{}{
        "id":         status.ID,
        "status":     status.Status, // pending, processing, completed, failed
        "attempts":   status.Attempts,
        "created_at": status.CreatedAt,
        "started_at": status.StartedAt,
        "finished_at": status.FinishedAt,
    }
}
```

### Queue Statistics

```go
func (c *AdminController) QueueStats(ctx types.Context) any {
    stats := c.queue.Stats()

    return map[string]interface{}{
        "queues": map[string]interface{}{
            "default": map[string]int64{
                "pending":    stats["default"].Pending,
                "processing": stats["default"].Processing,
                "completed":  stats["default"].Completed,
                "failed":     stats["default"].Failed,
            },
            "emails": map[string]int64{
                "pending":    stats["emails"].Pending,
                "processing": stats["emails"].Processing,
                "completed":  stats["emails"].Completed,
                "failed":     stats["emails"].Failed,
            },
        },
    }
}
```

## Job Middleware

Add middleware to jobs:

```go
type LoggingJobMiddleware struct{}

func (m *LoggingJobMiddleware) Handle(ctx queues.JobContext, next func() error) error {
    start := time.Now()
    log.Printf("Job %s starting", ctx.JobName())

    err := next()

    duration := time.Since(start)
    if err != nil {
        log.Printf("Job %s failed after %v: %v", ctx.JobName(), duration, err)
    } else {
        log.Printf("Job %s completed in %v", ctx.JobName(), duration)
    }

    return err
}
```

## Best Practices

1. **Keep jobs small** - Do one thing per job
2. **Make jobs idempotent** - Safe to run multiple times
3. **Set appropriate timeouts** - Prevent hung jobs
4. **Use retries** for transient failures
5. **Monitor queue depth** - Alert on backlogs
6. **Separate queues** by priority/type
7. **Log job activity** for debugging

## Common Use Cases

### Email Sending

```go
type WelcomeEmailJob struct {
    UserID string `json:"user_id"`
}

func (j *WelcomeEmailJob) Handle(ctx queues.JobContext) error {
    userService := ctx.Get("userService").(*UserService)
    user, _ := userService.GetByID(j.UserID)

    return sendWelcomeEmail(user.Email, user.Name)
}
```

### File Processing

```go
type ProcessUploadJob struct {
    FileID string `json:"file_id"`
}

func (j *ProcessUploadJob) Handle(ctx queues.JobContext) error {
    // Download file
    // Process (resize, convert, etc.)
    // Upload to storage
    // Update database
    return nil
}

func (j *ProcessUploadJob) Timeout() time.Duration {
    return 30 * time.Minute
}
```

### Report Generation

```go
type GenerateReportJob struct {
    ReportID string `json:"report_id"`
    UserID   string `json:"user_id"`
}

func (j *GenerateReportJob) Handle(ctx queues.JobContext) error {
    // Generate report
    // Save to storage
    // Notify user
    return nil
}
```

## Next Steps

- [Cron Jobs](cron.md) - Scheduled tasks
- [Caching](caching.md) - Performance optimization
- [Events](events.md) - Event-driven architecture
