# Events

Implement event-driven architecture in Goose applications.

## Overview

Events enable loose coupling between components:

- Publish events when things happen
- Subscribe to events to react
- Decouple services
- Enable extensibility

## Defining Events

### Basic Event

```go
package events

type UserRegistered struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Timestamp time.Time `json:"timestamp"`
}

func (e *UserRegistered) EventName() string {
    return "user.registered"
}
```

### Event with Metadata

```go
type OrderPlaced struct {
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Total     float64   `json:"total"`
    Items     []Item    `json:"items"`
    Timestamp time.Time `json:"timestamp"`
}

func (e *OrderPlaced) EventName() string {
    return "order.placed"
}

func (e *OrderPlaced) Metadata() map[string]interface{} {
    return map[string]interface{}{
        "priority": "high",
        "version":  "1.0",
    }
}
```

## Event Dispatcher

### Simple Dispatcher

```go
type EventDispatcher struct {
    handlers map[string][]EventHandler
    mutex    sync.RWMutex
}

type EventHandler func(event interface{}) error

func NewEventDispatcher() *EventDispatcher {
    return &EventDispatcher{
        handlers: make(map[string][]EventHandler),
    }
}

func (d *EventDispatcher) Subscribe(eventName string, handler EventHandler) {
    d.mutex.Lock()
    defer d.mutex.Unlock()

    d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *EventDispatcher) Dispatch(event Event) error {
    d.mutex.RLock()
    handlers := d.handlers[event.EventName()]
    d.mutex.RUnlock()

    for _, handler := range handlers {
        if err := handler(event); err != nil {
            return err
        }
    }

    return nil
}

func (d *EventDispatcher) DispatchAsync(event Event) {
    d.mutex.RLock()
    handlers := d.handlers[event.EventName()]
    d.mutex.RUnlock()

    for _, handler := range handlers {
        go handler(event)
    }
}
```

## Event Listeners

### Define Listener

```go
type SendWelcomeEmailListener struct {
    emailService *EmailService `inject:""`
}

func (l *SendWelcomeEmailListener) Handle(event interface{}) error {
    e := event.(*UserRegistered)
    return l.emailService.SendWelcomeEmail(e.Email, e.Name)
}

func (l *SendWelcomeEmailListener) ListensTo() []string {
    return []string{"user.registered"}
}
```

### Multiple Events

```go
type AuditLogListener struct {
    db *gorm.DB `inject:""`
}

func (l *AuditLogListener) Handle(event interface{}) error {
    // Log all events
    log := &AuditLog{
        EventType: event.(Event).EventName(),
        Payload:   toJSON(event),
        Timestamp: time.Now(),
    }
    return l.db.Create(log).Error
}

func (l *AuditLogListener) ListensTo() []string {
    return []string{
        "user.registered",
        "user.updated",
        "user.deleted",
        "order.placed",
        "order.shipped",
    }
}
```

## Publishing Events

### In Service

```go
type UserService struct {
    db         *gorm.DB         `inject:""`
    dispatcher *EventDispatcher `inject:""`
}

func (s *UserService) Register(dto RegisterDTO) (*User, error) {
    // Create user
    user := &User{
        ID:    uuid.New().String(),
        Email: dto.Email,
        Name:  dto.Name,
    }

    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }

    // Dispatch event
    s.dispatcher.DispatchAsync(&UserRegistered{
        UserID:    user.ID,
        Email:     user.Email,
        Name:      user.Name,
        Timestamp: time.Now(),
    })

    return user, nil
}
```

### In Controller

```go
type OrderController struct {
    orderService *OrderService    `inject:""`
    dispatcher   *EventDispatcher `inject:""`
}

func (c *OrderController) Create(ctx types.Context) any {
    var dto CreateOrderDTO
    ctx.Bind(&dto)

    order, err := c.orderService.Create(dto)
    if err != nil {
        return ctx.Status(500).JSON(map[string]string{"error": err.Error()})
    }

    // Dispatch event asynchronously
    c.dispatcher.DispatchAsync(&OrderPlaced{
        OrderID:   order.ID,
        UserID:    order.UserID,
        Total:     order.Total,
        Items:     order.Items,
        Timestamp: time.Now(),
    })

    return ctx.Status(201).JSON(order)
}
```

## Registering Listeners

### In Module

```go
func (m *AppModule) Declarations() []any {
    return []any{
        // Event listeners
        &SendWelcomeEmailListener{},
        &CreateUserProfileListener{},
        &AuditLogListener{},
        &NotifyAdminsListener{},
    }
}
```

### Manual Registration

```go
func setupEventListeners(dispatcher *EventDispatcher, services *Services) {
    // User events
    dispatcher.Subscribe("user.registered", func(e interface{}) error {
        event := e.(*UserRegistered)
        return services.Email.SendWelcomeEmail(event.Email, event.Name)
    })

    dispatcher.Subscribe("user.registered", func(e interface{}) error {
        event := e.(*UserRegistered)
        return services.Profile.CreateDefaultProfile(event.UserID)
    })

    // Order events
    dispatcher.Subscribe("order.placed", func(e interface{}) error {
        event := e.(*OrderPlaced)
        return services.Notification.NotifyOrderPlaced(event.OrderID)
    })
}
```

## Event Queue Integration

### Queue Events for Background Processing

```go
type QueuedEventDispatcher struct {
    queue *queues.Client `inject:""`
}

func (d *QueuedEventDispatcher) Dispatch(event Event) error {
    return d.queue.Dispatch(&ProcessEventJob{
        EventName: event.EventName(),
        Payload:   toJSON(event),
    })
}

type ProcessEventJob struct {
    EventName string `json:"event_name"`
    Payload   string `json:"payload"`
}

func (j *ProcessEventJob) Handle(ctx queues.JobContext) error {
    dispatcher := ctx.Get("eventDispatcher").(*EventDispatcher)

    // Reconstruct event
    event := reconstructEvent(j.EventName, j.Payload)

    // Process synchronously in background
    return dispatcher.ProcessSync(event)
}
```

## Event Patterns

### Domain Events

```go
// Base domain event
type DomainEvent struct {
    ID          string    `json:"id"`
    OccurredAt  time.Time `json:"occurred_at"`
    AggregateID string    `json:"aggregate_id"`
}

// Specific domain event
type ProductPriceChanged struct {
    DomainEvent
    ProductID string  `json:"product_id"`
    OldPrice  float64 `json:"old_price"`
    NewPrice  float64 `json:"new_price"`
}

func (e *ProductPriceChanged) EventName() string {
    return "product.price_changed"
}
```

### Event Sourcing

```go
type EventStore struct {
    db *gorm.DB `inject:""`
}

type StoredEvent struct {
    ID            string    `gorm:"primaryKey"`
    AggregateType string    `gorm:"index"`
    AggregateID   string    `gorm:"index"`
    EventType     string    `gorm:"index"`
    Payload       string    `gorm:"type:text"`
    OccurredAt    time.Time `gorm:"index"`
    Version       int
}

func (s *EventStore) Store(event Event) error {
    stored := &StoredEvent{
        ID:            uuid.New().String(),
        AggregateType: getAggregateType(event),
        AggregateID:   getAggregateID(event),
        EventType:     event.EventName(),
        Payload:       toJSON(event),
        OccurredAt:    time.Now(),
    }
    return s.db.Create(stored).Error
}

func (s *EventStore) GetEvents(aggregateID string) ([]StoredEvent, error) {
    var events []StoredEvent
    err := s.db.Where("aggregate_id = ?", aggregateID).
        Order("occurred_at ASC").
        Find(&events).Error
    return events, err
}
```

### Saga Pattern

```go
type OrderSaga struct {
    step int
    orderID string
    dispatcher *EventDispatcher
}

func (s *OrderSaga) Start(orderID string) {
    s.orderID = orderID
    s.dispatcher.Subscribe("payment.completed", s.onPaymentCompleted)
    s.dispatcher.Subscribe("payment.failed", s.onPaymentFailed)
    s.dispatcher.Subscribe("inventory.reserved", s.onInventoryReserved)
    s.dispatcher.Subscribe("inventory.failed", s.onInventoryFailed)
}

func (s *OrderSaga) onPaymentCompleted(e interface{}) error {
    event := e.(*PaymentCompleted)
    if event.OrderID != s.orderID {
        return nil
    }

    s.step++
    // Trigger next step
    return s.dispatcher.Dispatch(&ReserveInventory{OrderID: s.orderID})
}

func (s *OrderSaga) onPaymentFailed(e interface{}) error {
    event := e.(*PaymentFailed)
    if event.OrderID != s.orderID {
        return nil
    }

    // Compensate
    return s.dispatcher.Dispatch(&CancelOrder{
        OrderID: s.orderID,
        Reason:  "Payment failed",
    })
}
```

## Webhooks

Send events to external systems:

```go
type WebhookListener struct {
    httpClient *http.Client
    db         *gorm.DB `inject:""`
}

func (l *WebhookListener) Handle(event interface{}) error {
    e := event.(Event)

    // Get registered webhooks for this event
    var webhooks []Webhook
    l.db.Where("event_type = ?", e.EventName()).Find(&webhooks)

    // Send to each webhook
    for _, webhook := range webhooks {
        go l.sendWebhook(webhook, event)
    }

    return nil
}

func (l *WebhookListener) sendWebhook(webhook Webhook, event interface{}) {
    payload, _ := json.Marshal(event)

    req, _ := http.NewRequest("POST", webhook.URL, bytes.NewBuffer(payload))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-Webhook-Secret", webhook.Secret)

    resp, err := l.httpClient.Do(req)
    if err != nil {
        log.Printf("Webhook failed: %v", err)
        return
    }
    defer resp.Body.Close()
}
```

## Best Practices

1. **Name events in past tense** - `UserRegistered`, not `RegisterUser`
2. **Include timestamp** in all events
3. **Keep events immutable** - Don't modify after creation
4. **Make handlers idempotent** - Safe to process twice
5. **Use async dispatch** for non-critical handlers
6. **Log event processing** for debugging
7. **Handle failures gracefully** - Don't break main flow
8. **Version your events** for backward compatibility

## Common Events

```go
// User events
"user.registered"
"user.updated"
"user.deleted"
"user.logged_in"
"user.password_changed"

// Order events
"order.placed"
"order.paid"
"order.shipped"
"order.delivered"
"order.cancelled"

// Product events
"product.created"
"product.updated"
"product.deleted"
"product.out_of_stock"

// System events
"system.startup"
"system.shutdown"
"system.error"
```

## Next Steps

- [Queues](queues.md) - Background processing
- [Cron](cron.md) - Scheduled tasks
- [Services](../building-blocks/services.md) - Service layer
