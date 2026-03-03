# SQL Database

The SQL module provides GORM integration for relational databases.

## Quick Start

```go
import (
    "github.com/awesome-goose/goose/modules/sql"
)

// Configure the module
sqlModule := sql.NewModule(
    sql.WithDialect("sqlite"),
    sql.WithName("app.db"),
)

// Include in application
stop, err := goose.Start(goose.API(platform, module, []types.Module{
    sqlModule,
}))
```

## Configuration Options

### SQLite

```go
sqlModule := sql.NewModule(
    sql.WithDialect("sqlite"),
    sql.WithName("app.db"),
)
```

### PostgreSQL

```go
sqlModule := sql.NewModule(
    sql.WithDialect("postgres"),
    sql.WithHost("localhost"),
    sql.WithPort(5432),
    sql.WithName("myapp"),
    sql.WithUser("postgres"),
    sql.WithPass("secret"),
    sql.WithSSLMode("disable"),
)
```

### MySQL

```go
sqlModule := sql.NewModule(
    sql.WithDialect("mysql"),
    sql.WithHost("localhost"),
    sql.WithPort(3306),
    sql.WithName("myapp"),
    sql.WithUser("root"),
    sql.WithPass("secret"),
)
```

### Environment Configuration

```go
sqlModule := sql.NewModule(
    sql.WithDialect(env.String("DB_DIALECT", "sqlite")),
    sql.WithHost(env.String("DB_HOST", "localhost")),
    sql.WithPort(env.Int("DB_PORT", 5432)),
    sql.WithName(env.String("DB_NAME", "app.db")),
    sql.WithUser(env.String("DB_USER", "")),
    sql.WithPass(env.String("DB_PASS", "")),
)
```

### All Available Options

```go
type Config struct {
    Dialect    string       // "sqlite", "postgres", "mysql"
    Host       string       // Database host
    Port       int          // Database port
    User       string       // Database user
    Pass       string       // Database password
    Name       string       // Database name
    Sync       bool         // Auto-migrate tables
    Log        bool         // Enable query logging
    SSLMode    string       // SSL mode (postgres)
    Schema     string       // Schema name (postgres)
    TimeZone   string       // Timezone
    Seeders    []Seeder     // Database seeders
    Migrations []Migration  // Database migrations
}
```

**.env file:**

```env
DB_DIALECT=postgres
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASS=secret
DB_NAME=myapp
```

## Defining Entities

Entities are Go structs with GORM tags:

```go
package entities

import "time"

type User struct {
    ID        string     `json:"id" gorm:"primaryKey;type:uuid"`
    Email     string     `json:"email" gorm:"uniqueIndex;size:255"`
    Name      string     `json:"name" gorm:"size:100"`
    Age       int        `json:"age"`
    Active    bool       `json:"active" gorm:"default:true"`
    CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt *time.Time `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}
```

### Common GORM Tags

| Tag              | Description        |
| ---------------- | ------------------ |
| `primaryKey`     | Primary key field  |
| `uniqueIndex`    | Unique index       |
| `index`          | Regular index      |
| `size:255`       | Column size        |
| `type:uuid`      | Database type      |
| `default:value`  | Default value      |
| `autoCreateTime` | Auto-set on create |
| `autoUpdateTime` | Auto-set on update |
| `->`             | Read-only          |
| `-`              | Skip field         |

## Injecting the Database

Services receive the database via dependency injection:

```go
import "gorm.io/gorm"

type UserService struct {
    db *gorm.DB `inject:""`
}
```

## CRUD Operations

### Create

```go
func (s *UserService) Create(dto CreateUserDTO) (*User, error) {
    user := &User{
        ID:    uuid.New().String(),
        Email: dto.Email,
        Name:  dto.Name,
    }

    result := s.db.Create(user)
    if result.Error != nil {
        return nil, result.Error
    }

    return user, nil
}
```

### Read (Single)

```go
func (s *UserService) GetByID(id string) (*User, error) {
    var user User
    result := s.db.First(&user, "id = ?", id)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}

func (s *UserService) GetByEmail(email string) (*User, error) {
    var user User
    result := s.db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
```

### Read (Multiple)

```go
func (s *UserService) GetAll() []User {
    var users []User
    s.db.Find(&users)
    return users
}

func (s *UserService) GetActive() []User {
    var users []User
    s.db.Where("active = ?", true).Find(&users)
    return users
}

func (s *UserService) GetPaginated(page, limit int) []User {
    var users []User
    offset := (page - 1) * limit
    s.db.Offset(offset).Limit(limit).Find(&users)
    return users
}
```

### Update

```go
func (s *UserService) Update(id string, dto UpdateUserDTO) (*User, error) {
    var user User
    if err := s.db.First(&user, "id = ?", id).Error; err != nil {
        return nil, err
    }

    // Update fields
    user.Name = dto.Name
    user.Email = dto.Email

    if err := s.db.Save(&user).Error; err != nil {
        return nil, err
    }

    return &user, nil
}

// Partial update
func (s *UserService) UpdateName(id, name string) error {
    return s.db.Model(&User{}).Where("id = ?", id).Update("name", name).Error
}

// Update multiple fields
func (s *UserService) UpdateFields(id string, updates map[string]interface{}) error {
    return s.db.Model(&User{}).Where("id = ?", id).Updates(updates).Error
}
```

### Delete

```go
func (s *UserService) Delete(id string) error {
    return s.db.Delete(&User{}, "id = ?", id).Error
}

// Soft delete (requires DeletedAt field)
func (s *UserService) SoftDelete(id string) error {
    return s.db.Delete(&User{}, "id = ?", id).Error
}

// Hard delete
func (s *UserService) HardDelete(id string) error {
    return s.db.Unscoped().Delete(&User{}, "id = ?", id).Error
}
```

## Query Building

### Where Clauses

```go
// Simple
s.db.Where("name = ?", "John").Find(&users)

// Multiple conditions
s.db.Where("name = ? AND age > ?", "John", 18).Find(&users)

// IN clause
s.db.Where("id IN ?", []string{"1", "2", "3"}).Find(&users)

// LIKE
s.db.Where("name LIKE ?", "%john%").Find(&users)

// Struct condition
s.db.Where(&User{Name: "John", Active: true}).Find(&users)

// Map condition
s.db.Where(map[string]interface{}{"name": "John", "active": true}).Find(&users)
```

### Ordering and Limiting

```go
// Order
s.db.Order("created_at DESC").Find(&users)

// Multiple order
s.db.Order("name ASC, age DESC").Find(&users)

// Limit and offset
s.db.Limit(10).Offset(0).Find(&users)
```

### Selecting Fields

```go
// Select specific fields
s.db.Select("id", "name", "email").Find(&users)

// Exclude fields
s.db.Omit("password").Find(&users)
```

### Aggregations

```go
// Count
var count int64
s.db.Model(&User{}).Count(&count)

// Count with condition
s.db.Model(&User{}).Where("active = ?", true).Count(&count)
```

## Relationships

### One-to-Many

```go
type User struct {
    ID    string `gorm:"primaryKey"`
    Name  string
    Posts []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
    ID      string `gorm:"primaryKey"`
    Title   string
    UserID  string
    User    User `gorm:"foreignKey:UserID"`
}
```

### Preloading

```go
// Preload related
func (s *UserService) GetWithPosts(id string) (*User, error) {
    var user User
    result := s.db.Preload("Posts").First(&user, "id = ?", id)
    return &user, result.Error
}

// Nested preload
s.db.Preload("Posts.Comments").First(&user, "id = ?", id)

// Conditional preload
s.db.Preload("Posts", "published = ?", true).First(&user, "id = ?", id)
```

## Transactions

```go
func (s *OrderService) CreateOrder(dto CreateOrderDTO) (*Order, error) {
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // Create order
    order := &Order{ID: uuid.New().String(), UserID: dto.UserID}
    if err := tx.Create(order).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // Create line items
    for _, item := range dto.Items {
        lineItem := &LineItem{
            OrderID:   order.ID,
            ProductID: item.ProductID,
            Quantity:  item.Quantity,
        }
        if err := tx.Create(lineItem).Error; err != nil {
            tx.Rollback()
            return nil, err
        }
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return order, nil
}
```

## Migrations

Auto-migrate on startup:

```go
type AppService struct {
    db *gorm.DB `inject:""`
}

func (s *AppService) OnStart() {
    err := s.db.AutoMigrate(
        &User{},
        &Post{},
        &Comment{},
    )
    if err != nil {
        panic(err)
    }
}
```

## Raw SQL

```go
// Raw query
var users []User
s.db.Raw("SELECT * FROM users WHERE age > ?", 18).Scan(&users)

// Raw execute
s.db.Exec("UPDATE users SET active = ? WHERE last_login < ?", false, time.Now().AddDate(0, -6, 0))
```

## Best Practices

1. **Use transactions** for multiple related operations
2. **Index frequently queried columns**
3. **Use soft deletes** for audit trails
4. **Preload relationships** to avoid N+1 queries
5. **Use pagination** for large datasets
6. **Handle errors** from all database operations

## Next Steps

- [Entities](entities.md) - Entity definitions
- [Migrations](migrations.md) - Database migrations
- [KV Store](kv.md) - Key-value storage
