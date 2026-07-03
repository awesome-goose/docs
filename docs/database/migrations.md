# Database Migrations

Manage database schema changes with GORM migrations.

## Auto Migration

The simplest approach for development:

```go
type AppService struct {
    db *gorm.DB `inject:""`
}

// Boot runs once at startup (implements types.Bootable).
func (s *AppService) Boot(k types.Kernel) error {
    return s.db.AutoMigrate(
        &entities.User{},
        &entities.Post{},
        &entities.Comment{},
    )
}
```

### What AutoMigrate Does

- Creates tables if they don't exist
- Adds missing columns
- Adds missing indexes
- **Does NOT** drop columns
- **Does NOT** drop indexes
- **Does NOT** modify existing column types

## Migration on Startup

Configure migrations in your module:

```go
// app/app.module.go
package app

import (
    "myapp/app/entities"
    "github.com/awesome-goose/goose/types"
)

type AppModule struct{}

func (m *AppModule) Declarations() []any {
    return []any{
        &AppService{},
        &MigrationService{},
    }
}
```

```go
// app/migration.service.go
package app

import (
    "myapp/app/entities"
    "gorm.io/gorm"
)

type MigrationService struct {
    db *gorm.DB `inject:""`
}

func (s *MigrationService) Boot(k types.Kernel) error {
    s.runMigrations()
    return nil
}

func (s *MigrationService) runMigrations() {
    // Register all entities
    entities := []interface{}{
        &entities.User{},
        &entities.Profile{},
        &entities.Post{},
        &entities.Comment{},
        &entities.Tag{},
    }

    for _, entity := range entities {
        if err := s.db.AutoMigrate(entity); err != nil {
            panic(err)
        }
    }
}
```

## Conditional Migrations

Run migrations only in development:

```go
import "github.com/awesome-goose/goose/types"

type MigrationService struct {
    env types.Env `inject:""`
    // ...
}

func (s *MigrationService) Boot(k types.Kernel) error {
    if s.env.GetWithDefault("APP_ENV", "development") == "development" {
        s.runMigrations()
    }
    return nil
}
```

## Manual Migrations

For production, use explicit SQL:

```go
type MigrationService struct {
    db *gorm.DB `inject:""`
}

func (s *MigrationService) Boot(k types.Kernel) error {
    s.runMigrations()
    return nil
}

func (s *MigrationService) runMigrations() {
    migrations := []func(*gorm.DB) error{
        s.migration001CreateUsers,
        s.migration002AddEmailIndex,
        s.migration003CreatePosts,
    }

    for _, migrate := range migrations {
        if err := migrate(s.db); err != nil {
            panic(err)
        }
    }
}

func (s *MigrationService) migration001CreateUsers(db *gorm.DB) error {
    return db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id UUID PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            name VARCHAR(100),
            created_at TIMESTAMP DEFAULT NOW()
        )
    `).Error
}

func (s *MigrationService) migration002AddEmailIndex(db *gorm.DB) error {
    return db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)
    `).Error
}

func (s *MigrationService) migration003CreatePosts(db *gorm.DB) error {
    return db.Exec(`
        CREATE TABLE IF NOT EXISTS posts (
            id UUID PRIMARY KEY,
            title VARCHAR(255) NOT NULL,
            content TEXT,
            author_id UUID REFERENCES users(id),
            created_at TIMESTAMP DEFAULT NOW()
        )
    `).Error
}
```

## Migration Tracking

Track which migrations have run:

```go
type Migration struct {
    ID        string    `gorm:"primaryKey"`
    Name      string    `gorm:"uniqueIndex"`
    AppliedAt time.Time `gorm:"autoCreateTime"`
}

type MigrationService struct {
    db *gorm.DB `inject:""`
}

func (s *MigrationService) Boot(k types.Kernel) error {
    // Create migrations table
    s.db.AutoMigrate(&Migration{})

    // Run pending migrations
    s.runPending()
    return nil
}

func (s *MigrationService) runPending() {
    migrations := map[string]func(*gorm.DB) error{
        "001_create_users":     s.createUsers,
        "002_add_email_index":  s.addEmailIndex,
        "003_create_posts":     s.createPosts,
    }

    for name, migrate := range migrations {
        if !s.hasRun(name) {
            if err := migrate(s.db); err != nil {
                panic(fmt.Sprintf("Migration %s failed: %v", name, err))
            }
            s.markAsRun(name)
        }
    }
}

func (s *MigrationService) hasRun(name string) bool {
    var count int64
    s.db.Model(&Migration{}).Where("name = ?", name).Count(&count)
    return count > 0
}

func (s *MigrationService) markAsRun(name string) {
    s.db.Create(&Migration{
        ID:   uuid.New().String(),
        Name: name,
    })
}
```

## Schema Modifications

### Add Column

```go
func (s *MigrationService) addColumn(db *gorm.DB) error {
    // Check if column exists
    if db.Migrator().HasColumn(&User{}, "phone") {
        return nil
    }

    return db.Migrator().AddColumn(&User{}, "Phone")
}
```

### Rename Column

```go
func (s *MigrationService) renameColumn(db *gorm.DB) error {
    return db.Migrator().RenameColumn(&User{}, "name", "full_name")
}
```

### Drop Column

```go
func (s *MigrationService) dropColumn(db *gorm.DB) error {
    if !db.Migrator().HasColumn(&User{}, "deprecated_field") {
        return nil
    }

    return db.Migrator().DropColumn(&User{}, "deprecated_field")
}
```

### Modify Column

```go
func (s *MigrationService) modifyColumn(db *gorm.DB) error {
    // Change column type
    return db.Exec("ALTER TABLE users ALTER COLUMN name TYPE VARCHAR(200)").Error
}
```

### Add Index

```go
func (s *MigrationService) addIndex(db *gorm.DB) error {
    if db.Migrator().HasIndex(&User{}, "idx_users_email") {
        return nil
    }

    return db.Migrator().CreateIndex(&User{}, "idx_users_email")
}
```

### Drop Index

```go
func (s *MigrationService) dropIndex(db *gorm.DB) error {
    if !db.Migrator().HasIndex(&User{}, "idx_old_index") {
        return nil
    }

    return db.Migrator().DropIndex(&User{}, "idx_old_index")
}
```

### Add Foreign Key

```go
func (s *MigrationService) addForeignKey(db *gorm.DB) error {
    return db.Migrator().CreateConstraint(&Post{}, "Author")
}
```

## Data Migrations

Migrate data between schemas:

```go
func (s *MigrationService) migrateData(db *gorm.DB) error {
    // Start transaction
    tx := db.Begin()

    // Migrate data
    if err := tx.Exec(`
        UPDATE users
        SET full_name = CONCAT(first_name, ' ', last_name)
        WHERE full_name IS NULL
    `).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
```

## Rollback Support

Implement reversible migrations:

```go
type ReversibleMigration struct {
    Name string
    Up   func(*gorm.DB) error
    Down func(*gorm.DB) error
}

var migrations = []ReversibleMigration{
    {
        Name: "001_create_users",
        Up: func(db *gorm.DB) error {
            return db.AutoMigrate(&User{})
        },
        Down: func(db *gorm.DB) error {
            return db.Migrator().DropTable(&User{})
        },
    },
}

func (s *MigrationService) Rollback(name string) error {
    for _, m := range migrations {
        if m.Name == name {
            return m.Down(s.db)
        }
    }
    return fmt.Errorf("migration not found: %s", name)
}
```

## CLI Migrations

Create a CLI command for migrations:

```go
// app/cli/migrations.controller.go
type MigrationsController struct {
    db *gorm.DB `inject:""`
}

type EmptyDto struct{}

func (c *MigrationsController) Migrate(dto *EmptyDto) types.Output {
    // Run migrations
    return output.ConsoleSuccess("Migrations complete")
}

func (c *MigrationsController) Status(dto *EmptyDto) types.Output {
    // Show migration status
    return output.Line("Pending: 2, Applied: 5")
}

func (c *MigrationsController) Rollback(dto *EmptyDto) types.Output {
    // Rollback last migration
    return output.ConsoleSuccess("Rolled back: 007_add_indexes")
}
```

Register the commands with the CLI router:

```go
// app/cli/migrations.routes.go
var ROUTES = router.ForRoutes(
    router.Cli("/migrate", []any{MigrationsController{}, "Migrate"}),
    router.Cli("/migrate/status", []any{MigrationsController{}, "Status"}),
    router.Cli("/migrate/rollback", []any{MigrationsController{}, "Rollback"}),
)
```

## Best Practices

1. **Track migrations** in a database table
2. **Test migrations** on a copy of production data
3. **Make migrations idempotent** (safe to run twice)
4. **Backup before migrating** in production
5. **Use transactions** for data migrations
6. **Keep migrations small** and focused
7. **Never edit** applied migrations
8. **Document breaking changes**

## Production Checklist

- [ ] Backup database before migration
- [ ] Test migration on staging
- [ ] Plan for rollback
- [ ] Schedule maintenance window
- [ ] Monitor during migration
- [ ] Verify data integrity after

## Next Steps

- [SQL Database](sql.md) - Database operations
- [Entities](entities.md) - Entity definitions
- [Deployment](../deployment/overview.md) - Production deployment
