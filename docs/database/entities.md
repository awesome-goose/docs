# Database Entities

Entities define the structure of your database tables using Go structs with GORM tags.

## Basic Entity

```go
package entities

import "time"

type User struct {
    ID        string     `json:"id" gorm:"primaryKey"`
    Email     string     `json:"email" gorm:"uniqueIndex"`
    Name      string     `json:"name"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}
```

## GORM Tag Reference

### Primary Key

```go
type Entity struct {
    // UUID primary key
    ID string `gorm:"primaryKey;type:uuid"`

    // Auto-increment integer
    ID uint `gorm:"primaryKey;autoIncrement"`
}
```

### Column Types

```go
type Entity struct {
    // String with size
    Name string `gorm:"size:100"`

    // Text (unlimited)
    Description string `gorm:"type:text"`

    // Boolean with default
    Active bool `gorm:"default:true"`

    // Integer
    Age int `gorm:"type:int"`

    // Float
    Price float64 `gorm:"type:decimal(10,2)"`

    // JSON field
    Metadata datatypes.JSON `gorm:"type:json"`
}
```

### Indexes

```go
type Entity struct {
    // Unique index
    Email string `gorm:"uniqueIndex"`

    // Regular index
    Status string `gorm:"index"`

    // Named index
    Name string `gorm:"index:idx_name"`

    // Composite index
    FirstName string `gorm:"index:idx_full_name,priority:1"`
    LastName  string `gorm:"index:idx_full_name,priority:2"`
}
```

### Constraints

```go
type Entity struct {
    // Not null (default)
    Name string

    // Nullable
    Phone *string `gorm:"default:null"`

    // Default value
    Status string `gorm:"default:'pending'"`

    // Check constraint
    Age int `gorm:"check:age >= 0"`
}
```

### Timestamps

```go
type Entity struct {
    // Auto-set on create
    CreatedAt time.Time `gorm:"autoCreateTime"`

    // Auto-set on update
    UpdatedAt time.Time `gorm:"autoUpdateTime"`

    // Unix timestamp
    CreatedAt int64 `gorm:"autoCreateTime:milli"`

    // Soft delete
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

## Common Entity Patterns

### Base Entity

Create a reusable base:

```go
package entities

import (
    "time"
    "gorm.io/gorm"
)

type BaseEntity struct {
    ID        string         `json:"id" gorm:"primaryKey;type:uuid"`
    CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// Generate ID before create
func (e *BaseEntity) BeforeCreate(tx *gorm.DB) error {
    if e.ID == "" {
        e.ID = uuid.New().String()
    }
    return nil
}
```

Use in entities:

```go
type User struct {
    BaseEntity
    Email string `json:"email" gorm:"uniqueIndex"`
    Name  string `json:"name"`
}

type Product struct {
    BaseEntity
    Name  string  `json:"name"`
    Price float64 `json:"price"`
}
```

### Soft Delete

Enable soft delete:

```go
import "gorm.io/gorm"

type User struct {
    ID        string         `gorm:"primaryKey"`
    Email     string         `gorm:"uniqueIndex"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

// Soft delete
db.Delete(&user)  // Sets DeletedAt

// Query (excludes soft deleted)
db.Find(&users)   // WHERE deleted_at IS NULL

// Include soft deleted
db.Unscoped().Find(&users)

// Hard delete
db.Unscoped().Delete(&user)
```

## Relationships

### One-to-One

```go
type User struct {
    ID      string `gorm:"primaryKey"`
    Profile Profile
}

type Profile struct {
    ID     string `gorm:"primaryKey"`
    UserID string `gorm:"uniqueIndex"`
    Bio    string
}
```

### One-to-Many

```go
type User struct {
    ID    string `gorm:"primaryKey"`
    Posts []Post `gorm:"foreignKey:AuthorID"`
}

type Post struct {
    ID       string `gorm:"primaryKey"`
    Title    string
    AuthorID string
    Author   User `gorm:"foreignKey:AuthorID"`
}
```

### Many-to-Many

```go
type User struct {
    ID    string `gorm:"primaryKey"`
    Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    ID    string `gorm:"primaryKey"`
    Name  string
    Users []User `gorm:"many2many:user_roles;"`
}
```

### Self-Referential

```go
type Category struct {
    ID       string     `gorm:"primaryKey"`
    Name     string
    ParentID *string
    Parent   *Category  `gorm:"foreignKey:ParentID"`
    Children []Category `gorm:"foreignKey:ParentID"`
}
```

### Polymorphic

```go
type Comment struct {
    ID            string `gorm:"primaryKey"`
    Content       string
    CommentableID string
    CommentableType string
}

type Post struct {
    ID       string    `gorm:"primaryKey"`
    Comments []Comment `gorm:"polymorphic:Commentable;"`
}

type Video struct {
    ID       string    `gorm:"primaryKey"`
    Comments []Comment `gorm:"polymorphic:Commentable;"`
}
```

## Hooks

Lifecycle callbacks:

```go
type User struct {
    ID       string `gorm:"primaryKey"`
    Email    string
    Password string
}

// Before create
func (u *User) BeforeCreate(tx *gorm.DB) error {
    u.ID = uuid.New().String()
    return nil
}

// Before save (create and update)
func (u *User) BeforeSave(tx *gorm.DB) error {
    if u.Password != "" {
        hashed, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
        u.Password = string(hashed)
    }
    return nil
}

// After create
func (u *User) AfterCreate(tx *gorm.DB) error {
    // Send welcome email
    return nil
}

// Before delete
func (u *User) BeforeDelete(tx *gorm.DB) error {
    // Cleanup related data
    return nil
}
```

## JSON Serialization

Control JSON output:

```go
type User struct {
    ID       string `json:"id" gorm:"primaryKey"`
    Email    string `json:"email"`
    Password string `json:"-"`  // Never serialize
    Name     string `json:"name,omitempty"`
}
```

## Scopes

Reusable query scopes:

```go
type User struct {
    ID     string `gorm:"primaryKey"`
    Active bool
    Role   string
}

// Define scopes
func ActiveUsers(db *gorm.DB) *gorm.DB {
    return db.Where("active = ?", true)
}

func WithRole(role string) func(*gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        return db.Where("role = ?", role)
    }
}

// Use scopes
db.Scopes(ActiveUsers).Find(&users)
db.Scopes(ActiveUsers, WithRole("admin")).Find(&admins)
```

## Entity Organization

### Directory Structure

```
app/
├── entities/
│   ├── base.go      # BaseEntity
│   ├── user.go
│   ├── post.go
│   └── comment.go
└── users/
    ├── users.module.go
    ├── users.service.go
    └── users.controller.go
```

### Entity File

```go
// app/entities/user.go
package entities

import "time"

type User struct {
    ID        string     `json:"id" gorm:"primaryKey;type:uuid"`
    Email     string     `json:"email" gorm:"uniqueIndex;size:255"`
    Name      string     `json:"name" gorm:"size:100"`
    Password  string     `json:"-" gorm:"size:255"`
    Active    bool       `json:"active" gorm:"default:true"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`

    // Relationships
    Posts []Post `json:"posts,omitempty" gorm:"foreignKey:AuthorID"`
}

func (User) TableName() string {
    return "users"
}
```

## Best Practices

1. **Use UUID** for primary keys in distributed systems
2. **Add soft delete** for audit trails
3. **Create base entity** for common fields
4. **Use hooks** for automatic operations
5. **Define table names** explicitly
6. **Index foreign keys** for performance
7. **Validate data** before saving

## Next Steps

- [SQL Database](sql.md) - Database operations
- [Migrations](migrations.md) - Schema migrations
- [Validation](../building-blocks/validation.md) - Data validation
