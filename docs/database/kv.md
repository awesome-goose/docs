# Key-Value Store

The KV module provides a SQL-backed key-value store with TTL and groups
(namespaces). It is **not** a Redis client — it persists values through the
same GORM connection used by the SQL module, so all operations are durable.

## Quick Start

```go
import (
    "github.com/awesome-goose/goose/modules/kv"
    "github.com/awesome-goose/goose/modules/sql"
    "github.com/awesome-goose/goose/types"
)

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        sql.Root(&sql.Config{Dialect: "sqlite", Name: "app.db", Sync: true}),
        kv.NewModule(&kv.Config{
            Group:           "app",
            DefaultTTL:      0,
            CleanupInterval: 0,
        }, false),
    }
}
```

## Configuration

```go
type Config struct {
    Group           string        // Namespace (similar to a Redis DB number)
    DefaultTTL      time.Duration // Default TTL for new keys (0 = no expiration)
    CleanupInterval time.Duration // How often to sweep expired keys (0 = never)
}
```

Functional setters are available too — apply them to a `*kv.Config` before
passing it to `NewModule`:

```go
cfg := &kv.Config{}
kv.WithGroup("sessions")(cfg)
kv.WithDefaultTTL(24 * time.Hour)(cfg)
kv.WithCleanupInterval(10 * time.Minute)(cfg)
kv.NewModule(cfg, false)
```

> The KV module does **not** take host/port/password/driver options — it
> piggybacks on the SQL module's connection.

## Injecting the Store

Services receive a `*kv.KV` via dependency injection:

```go
import "github.com/awesome-goose/goose/modules/kv"

type CacheService struct {
    kv *kv.KV `inject:""`
}
```

## Basic Operations

### Set / Get / Delete

```go
// Set
err := s.kv.Set("user:123", user)            // uses DefaultTTL
err = s.kv.Set("user:123", user, time.Hour)  // explicit TTL

// Get
v, err := s.kv.Get("user:123") // (any, error); ErrKeyNotFound if missing

// Delete
n, err := s.kv.Del("user:123", "user:124")   // returns count deleted
```

### Conditional Set

```go
// Only writes if the key does not already exist.
created, err := s.kv.SetNX("lock:resource", "owner-id", 30*time.Second)
if created {
    // got the lock
}
```

### Set-And-Get

```go
// Atomically write the new value and return the previous one.
old, err := s.kv.GetSet("counter", "0", time.Hour)
```

## Counters

```go
n, err := s.kv.Incr("views:post:42")           // returns new value
n, err = s.kv.IncrBy("views:post:42", 5)
```

`Incr` and `IncrBy` are atomic in the SQL backend.

## Expiration

```go
ok, err := s.kv.Expire("session:abc", 24*time.Hour) // refresh TTL
ok, err = s.kv.Persist("session:abc")              // drop expiration
ttl, err := s.kv.TTL("session:abc")                // seconds remaining
```

## Existence / Listing

```go
count, err := s.kv.Exists("user:123", "user:124") // number of present keys
keys, err := s.kv.Keys("user:*")                  // pattern match
```

## Caching Pattern

```go
type ProductService struct {
    kv *kv.KV  `inject:""`
    db *sql.Db `inject:""`
}

func (s *ProductService) GetByID(id string) (*Product, error) {
    cacheKey := "product:" + id

    if v, err := s.kv.Get(cacheKey); err == nil {
        if p, ok := v.(*Product); ok {
            return p, nil
        }
    }

    var p Product
    if err := s.db.First(&p, "id = ?", id).Error; err != nil {
        return nil, err
    }

    _ = s.kv.Set(cacheKey, &p, time.Hour)
    return &p, nil
}
```

## Session Management

```go
type SessionService struct {
    kv *kv.KV `inject:""`
}

func (s *SessionService) Create(userID string) (string, error) {
    sid := uuid.New().String()
    return sid, s.kv.Set("session:"+sid, map[string]any{
        "user_id":    userID,
        "created_at": time.Now(),
    }, 24*time.Hour)
}

func (s *SessionService) Get(sid string) (any, error) {
    return s.kv.Get("session:" + sid)
}

func (s *SessionService) Destroy(sid string) error {
    _, err := s.kv.Del("session:" + sid)
    return err
}

func (s *SessionService) Refresh(sid string) error {
    _, err := s.kv.Expire("session:"+sid, 24*time.Hour)
    return err
}
```

## Rate Limiting

```go
type RateLimiter struct {
    kv *kv.KV `inject:""`
}

func (r *RateLimiter) Allow(key string, limit int, window time.Duration) bool {
    countKey := "ratelimit:" + key

    if _, err := r.kv.Get(countKey); err != nil {
        _ = r.kv.Set(countKey, "1", window)
        return true
    }

    n, _ := r.kv.Incr(countKey)
    return n <= int64(limit)
}
```

## What this module is NOT

- It is not a Redis client. There is no `HSet`/`HGet`, `LPush`/`RPop`, or
  `SAdd`/`SMembers`. Encode complex shapes as JSON and store them under a
  single key, or use the SQL module for relational data.
- It does not take `Host`/`Port`/`Password`/`Driver` options. The store
  inherits its backend from the SQL module.

## Best Practices

1. **Use namespaced keys** with consistent patterns (`user:123:profile`).
2. **Set a TTL** for cached data via `DefaultTTL` or per-call.
3. **Handle missing keys** by checking `err == nil`; the value is `any` and
   may need a type assertion.
4. **Avoid huge values** — JSON-encoded blobs are fine; megabytes are not.

## Next Steps

- [SQL Database](sql.md) - Relational database
- [Caching](../advanced/caching.md) - Advanced caching
