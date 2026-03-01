# Database Overview

Goose provides built-in modules for database operations with support for SQL databases, key-value stores, and document databases.

## Database Modules

| Module            | Description              | Use Case                  |
| ----------------- | ------------------------ | ------------------------- |
| [SQL](sql.md)     | GORM-based SQL databases | PostgreSQL, MySQL, SQLite |
| [KV](kv.md)       | Key-value storage        | Redis, caching, sessions  |
| [NoSQL](nosql.md) | Document databases       | MongoDB (coming soon)     |

## SQL Database

The SQL module provides GORM integration:

```go
import (
    "myapp/app"
    "github.com/awesome-goose/goose"
    "github.com/awesome-goose/goose/modules/sql"
)

func main() {
    sqlModule := sql.NewModule(
        sql.WithDriver("sqlite"),
        sql.WithDSN("file:app.db?cache=shared&mode=rwc"),
    )

    stop, err := goose.Start(goose.API(platform, module, []types.Module{
        sqlModule,
    }))
}
```

### Supported Databases

- SQLite
- PostgreSQL
- MySQL/MariaDB
- SQL Server

## Key-Value Store

The KV module provides Redis-compatible storage:

```go
import "github.com/awesome-goose/goose/modules/kv"

kvModule := kv.NewModule(
    kv.WithDriver("redis"),
    kv.WithHost("localhost"),
    kv.WithPort(6379),
)
```

### Use Cases

- Caching
- Session storage
- Rate limiting
- Temporary data

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

Services access databases through injection:

```go
type UserService struct {
    db *gorm.DB `inject:""`
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
    tx := s.db.Begin()

    // Create order
    order := &Order{...}
    if err := tx.Create(order).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // Create line items
    for _, item := range dto.Items {
        lineItem := &LineItem{OrderID: order.ID, ...}
        if err := tx.Create(lineItem).Error; err != nil {
            tx.Rollback()
            return nil, err
        }
    }

    tx.Commit()
    return order, nil
}
```

## Migrations

Auto-migrate entities on startup:

```go
func (s *AppService) OnStart() {
    s.db.AutoMigrate(
        &User{},
        &Product{},
        &Order{},
    )
}
```

## Configuration

Configure via environment:

```env
# SQL Database
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=postgres
DB_PASSWORD=secret

# Redis KV
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Quick Start

1. Add the SQL module:

```go
sqlModule := sql.NewModule(
    sql.WithDriver(env.String("DB_DRIVER", "sqlite")),
    sql.WithDSN(env.String("DB_DSN", "file:app.db")),
)
```

2. Define an entity:

```go
type Product struct {
    ID    string  `gorm:"primaryKey"`
    Name  string
    Price float64
}
```

3. Inject and use:

```go
type ProductService struct {
    db *gorm.DB `inject:""`
}

func (s *ProductService) Create(name string, price float64) *Product {
    product := &Product{
        ID:    uuid.New().String(),
        Name:  name,
        Price: price,
    }
    s.db.Create(product)
    return product
}
```

## Next Steps

- [SQL Module](sql.md) - Detailed SQL usage
- [KV Module](kv.md) - Key-value operations
- [Entities](entities.md) - Entity definitions
- [Migrations](migrations.md) - Database migrations
