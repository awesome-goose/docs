# Database Overview

Goose provides built-in modules for database operations with support for SQL
databases and key-value stores.

## Database Modules

| Module        | Description                | Use Case                  |
| ------------- | -------------------------- | ------------------------- |
| [SQL](sql.md) | GORM-based SQL databases   | PostgreSQL, MySQL, SQLite |
| [KV](kv.md)   | Persistent key-value store | Sessions, flags, configs  |

## SQL Database

The SQL module provides GORM integration. Register it as a root or child
module from your `Imports()`:

```go
import (
    "github.com/awesome-goose/goose/modules/sql"
    "github.com/awesome-goose/goose/types"
)

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        sql.Root(&sql.Config{
            Dialect: "sqlite",
            Name:    "app.db",
            Sync:    true,
        }),
    }
}
```

`sql.Root(cfg)` and `sql.Child(cfg)` are shortcuts for
`sql.NewModule(cfg, true)` and `sql.NewModule(cfg, false)`.

### Supported Dialects

- `sqlite`
- `postgres`
- `mysql`

## Key-Value Store

The KV module persists values through the SQL connection (it uses the same
GORM database) and adds TTL/group semantics:

```go
import "github.com/awesome-goose/goose/modules/kv"

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        sql.Root(&sql.Config{Dialect: "sqlite", Name: "app.db", Sync: true}),
        kv.NewModule(&kv.Config{
            Group:           "app",
            DefaultTTL:      0, // 0 = no expiration
            CleanupInterval: 0,
        }, false),
    }
}
```

### Use Cases

- Feature flags
- Long-lived counters
- Lightweight configuration
- One-time tokens (with `DefaultTTL`)

## Database Entities

Define entities with GORM tags:

```go
type User struct {
    ID        string     `json:"id" gorm:"primaryKey"`
    Email     string     `json:"email" gorm:"uniqueIndex"`
    Name      string     `json:"name"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}
```

## Repository Pattern

Services access the database through dependency injection. Inject the
`*sql.Db` (a thin wrapper around `*gorm.DB`) or the higher-level `*sql.Query`
helper:

```go
import (
    "github.com/awesome-goose/goose/modules/sql"
)

type UserService struct {
    db *sql.Db `inject:""`
}

func (s *UserService) GetAll() []User {
    var users []User
    s.db.Find(&users)
    return users
}
```

## Transactions

Use transactions for data integrity:

```go
func (s *OrderService) CreateOrder(dto CreateOrderDTO) (*Order, error) {
    return s.db.Transaction(func(tx *gorm.DB) error {
        order := &Order{}
        if err := tx.Create(order).Error; err != nil {
            return err
        }
        for _, item := range dto.Items {
            li := &LineItem{OrderID: order.ID}
            if err := tx.Create(li).Error; err != nil {
                return err
            }
        }
        return nil
    })
}
```

## Migrations

Set `Sync: true` on the SQL config to auto-migrate registered entities at
startup, or run migrations explicitly via the SQL module's `Migrations`
slice. See [Migrations](migrations.md) for the explicit form.

## Configuration

Configure via environment, injecting `types.Env`:

```env
# SQL Database
DB_DIALECT=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=postgres
DB_PASS=secret
```

```go
type AppModule struct {
    env types.Env `inject:""`
}

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        sql.Root(&sql.Config{
            Dialect: m.env.GetWithDefault("DB_DIALECT", "sqlite"),
            Host:    m.env.GetWithDefault("DB_HOST", "localhost"),
            Port:    m.env.GetInt("DB_PORT"),
            Name:    m.env.GetWithDefault("DB_NAME", "app.db"),
            User:    m.env.Get("DB_USER"),
            Pass:    m.env.Get("DB_PASS"),
        }),
    }
}
```

## Next Steps

- [SQL Module](sql.md) - Detailed SQL usage
- [KV Module](kv.md) - Key-value operations
- [Entities](entities.md) - Entity definitions
- [Migrations](migrations.md) - Database migrations
