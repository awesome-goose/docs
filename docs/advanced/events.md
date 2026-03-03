# Events

> **Note:** The Events module is not a built-in Goose feature. This documentation describes an event-driven pattern that you can implement manually in your application.

## Overview

Events enable loose coupling between components:

- Publish events when things happen
- Subscribe to events to react
- Decouple services
- Enable extensibility

## Manual Implementation

### Event Dispatcher

```go
package events

import "sync"

type Event interface {
    EventName() string
}

type EventHandler func(event interface{}) error

type Dispatcher struct {
    handlers map[string][]EventHandler
    mutex    sync.RWMutex
}

func NewDispatcher() *Dispatcher {
    return &Dispatcher{
        handlers: make(map[string][]EventHandler),
    }
}

func (d *Dispatcher) Subscribe(eventName string, handler EventHandler) {
    d.mutex.Lock()
    defer d.mutex.Unlock()
    d.handlers[eventName] = append(d.handlers[eventName], handler)
}

func (d *Dispatcher) Dispatch(event Event) error {
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

func (d *Dispatcher) DispatchAsync(event Event) {
    d.mutex.RLock()
    handlers := d.handlers[event.EventName()]
    d.mutex.RUnlock()

    for _, handler := range handlers {
        go handler(event)
    }
}
```

## Defining Events

```go
package events

import "time"

type UserRegistered struct {
    UserID    string    `json:"user_id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    Timestamp time.Time `json:"timestamp"`
}

func (e *UserRegistered) EventName() string {
    return "user.registered"
}

type OrderPlaced struct {
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Total     float64   `json:"total"`
    Timestamp time.Time `json:"timestamp"`
}

func (e *OrderPlaced) EventName() string {
    return "order.placed"
}
```

## Event Listeners

```go
type SendWelcomeEmailListener struct {
    emailService *EmailService
}

func (l *SendWelcomeEmailListener) Handle(event interface{}) error {
    e := event.(*UserRegistered)
    return l.emailService.SendWelcomeEmail(e.Email, e.Name)
}
```

## Publishing Events

### In Service

```go
type UserService struct {
    db         *gorm.DB    `inject:""`
    dispatcher *Dispatcher `inject:""`
}

func (s *UserService) Register(dto RegisterDTO) (*User, error) {
    user := &User{
        ID:    uuid.New().String(),
        Email: dto.Email,
        Name:  dto.Name,
    }

    if err := s.db.Create(user).Error; err != nil {
        return nil, err
    }

    // Dispatch event asynchronously
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
    orderService *OrderService `inject:""`
    dispatcher   *Dispatcher   `inject:""`
}

func (c *OrderController) Create(ctx types.Context) any {
    body, _ := ctx.Request().Body()
    var dto CreateOrderDTO
    json.Unmarshal(body, &dto)

    order, err := c.orderService.Create(dto)
    if err != nil {
        return map[string]any{"error": err.Error(), "_status": 500}
    }

    c.dispatcher.DispatchAsync(&OrderPlaced{
        OrderID:   order.ID,
        UserID:    order.UserID,
        Total:     order.Total,
        Timestamp: time.Now(),
    })

    return map[string]any{"data": order, "_status": 201}
}
```

## Registering the Dispatcher

Register the dispatcher as a singleton using the `types.Configurable` interface:

```go
func (m *AppModule) Configure(c types.Container) error {
    dispatcher := events.NewDispatcher()

    // Subscribe handlers
    dispatcher.Subscribe("user.registered", func(e interface{}) error {
        event := e.(*events.UserRegistered)
        // Handle event (e.g., send welcome email)
        return nil
    })

    return c.Register(func() *events.Dispatcher {
        return dispatcher
    }, "", true)
}
```

## Using with Queues Module

For more reliable event processing, use the queues module:

```go
import "github.com/awesome-goose/goose/modules/queues"

type ProcessEventJob struct {
    EventName string `json:"event_name"`
    Payload   string `json:"payload"`
}

// Queue event processing instead of handling immediately
func (s *UserService) Register(dto RegisterDTO) (*User, error) {
    // ... create user ...

    // Queue the event processing for reliability
    payload, _ := json.Marshal(&UserRegistered{
        UserID:    user.ID,
        Email:     user.Email,
        Name:      user.Name,
        Timestamp: time.Now(),
    })

    s.queue.Dispatch(&ProcessEventJob{
        EventName: "user.registered",
        Payload:   string(payload),
    })

    return user, nil
}
```

## Best Practices

1. **Name events in past tense** - `UserRegistered`, not `RegisterUser`
2. **Include timestamps** - Always include when the event occurred
3. **Keep events immutable** - Don't modify events after creation
4. **Make handlers idempotent** - Safe to process multiple times
5. **Use async dispatch** for non-critical handlers
6. **Handle failures gracefully** - Don't break main flow
7. **Log event processing** for debugging

## Common Event Types

```go
// User events
"user.registered"
"user.updated"
"user.deleted"
"user.logged_in"

// Order events
"order.placed"
"order.paid"
"order.shipped"
"order.cancelled"

// Product events
"product.created"
"product.updated"
"product.out_of_stock"
```

## Next Steps

- [Queues](queues.md) - Background job processing
- [Cron](cron.md) - Scheduled tasks
- [Services](../building-blocks/services.md) - Business logic layer
