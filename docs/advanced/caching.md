# Caching

Improve performance with the Goose `cache` module.

## Overview

The Cache module provides a SQL-backed caching interface with support for:

- TTL (Time-To-Live)
- Cache groups for namespacing
- Automatic cleanup of expired entries
- Generic `Remember[T]` / `GetAs[T]` helpers

The cache uses the same SQL connection as `sql.Root` — there is no Redis
client. If you need durable cross-process caching, point the SQL module at
Postgres or MySQL.

## Quick Start

```go
import (
    "github.com/awesome-goose/goose/modules/cache"
    "github.com/awesome-goose/goose/modules/sql"
    "github.com/awesome-goose/goose/types"
)

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        sql.Root(&sql.Config{Dialect: "sqlite", Name: "app.db", Sync: true}),
        cache.NewModule(&cache.Config{
            Group:           "app",
            DefaultTTL:      time.Hour,
            CleanupInterval: 10 * time.Minute,
        }, false),
    }
}
```

## Configuration

```go
type Config struct {
    Group           string        // Namespace for cache entries
    DefaultTTL      time.Duration // Default TTL (e.g. 5*time.Minute)
    CleanupInterval time.Duration // Sweep interval (0 = never)
}
```

Functional setters are applied to the config before passing it to
`NewModule`:

```go
cfg := &cache.Config{}
cache.WithGroup("products")(cfg)
cache.WithDefaultTTL(time.Hour)(cfg)
cache.WithCleanupInterval(10 * time.Minute)(cfg)
cache.NewModule(cfg, false)
```

To read the group from the environment, inject `types.Env` and read it in
the module's `Imports()`:

```go
type AppModule struct {
    env types.Env `inject:""`
}

func (m *AppModule) Imports() []types.Module {
    return []types.Module{
        cache.NewModule(&cache.Config{
            Group:      m.env.GetWithDefault("CACHE_GROUP", "default"),
            DefaultTTL: time.Hour,
        }, false),
    }
}
```

## Using the Cache

### Inject the cache

```go
import "github.com/awesome-goose/goose/modules/cache"

type ProductService struct {
    cache *cache.Cache `inject:""`
    db    *sql.Db      `inject:""`
}
```

### Set / Get / Delete / Has

```go
// Set with TTL (defaults to Config.DefaultTTL if omitted)
err := s.cache.Set("product:42", product, time.Hour)

// Get returns (any, error). Use cache.GetAs[T] for a typed result.
v, err := s.cache.Get("product:42")

// Delete returns the number of keys removed
n, err := s.cache.Delete("product:42", "product:43")

// Has returns (true, nil) if the key is present and not expired
ok, err := s.cache.Has("product:42")
```

### Typed Get

`cache.GetAs[T]` returns a value of the requested type or
`cache.ErrKeyNotFound`:

```go
product, err := cache.GetAs[*Product](s.cache, "product:42")
if err != nil {
    return nil, err
}
return product, nil
```

### Invalidate / Flush

```go
err := s.cache.Invalidate("product:42")   // single key
n, err := s.cache.InvalidateAll()         // every key in the group
n, err = s.cache.Flush()                  // alias for InvalidateAll
```

## Caching Patterns

### Cache-Aside (Lazy Loading)

```go
func (s *ProductService) GetByID(id string) (*Product, error) {
    cacheKey := "product:" + id

    if p, err := cache.GetAs[*Product](s.cache, cacheKey); err == nil {
        return p, nil
    }

    var p Product
    if err := s.db.First(&p, "id = ?", id).Error; err != nil {
        return nil, err
    }

    _ = s.cache.Set(cacheKey, &p, time.Hour)
    return &p, nil
}
```

### Remember

`cache.Remember[T]` is a generic helper that returns the cached value or
populates it via a callback:

```go
func (s *ProductService) GetByID(id string) (*Product, error) {
    return cache.Remember[*Product](
        s.cache,
        "product:"+id,
        func() (*Product, error) {
            var p Product
            err := s.db.First(&p, "id = ?", id).Error
            return &p, err
        },
        time.Hour, // optional; defaults to Config.DefaultTTL
    )
}
```

There is no `RememberForever`. For long-lived entries, pass a large TTL
explicitly (e.g. `100*365*24*time.Hour`) or set `Config.DefaultTTL` to zero
and rely on the absence of expiration.

### Write-Through

```go
func (s *ProductService) Update(id string, dto UpdateProductDTO) (*Product, error) {
    var p Product
    if err := s.db.First(&p, "id = ?", id).Error; err != nil {
        return nil, err
    }
    p.Name, p.Price = dto.Name, dto.Price
    if err := s.db.Save(&p).Error; err != nil {
        return nil, err
    }

    _ = s.cache.Set("product:"+id, &p, time.Hour)
    return &p, nil
}
```

### Invalidate Related Keys

```go
func (s *ProductService) Delete(id string) error {
    if err := s.db.Delete(&Product{}, "id = ?", id).Error; err != nil {
        return err
    }
    _, _ = s.cache.Delete("product:"+id, "products:list")
    return nil
}
```

## Group-Scoped Invalidation

Goose does not have a `Tags(...)` API. Use the `Group` config option (or a
per-call key prefix) to namespace related entries, and flush by creating a
dedicated cache module for that group:

```go
// In a module that holds short-lived "category" cache:
cache.NewModule(&cache.Config{
    Group:      "category-cache",
    DefaultTTL: 10 * time.Minute,
}, false)

// Flush only that group:
n, _ := s.categoryCache.InvalidateAll()
```

## Caching Lists / Paginated Results

```go
func (s *ProductService) GetPaginated(page, limit int) ([]Product, error) {
    key := fmt.Sprintf("products:page:%d:limit:%d", page, limit)
    return cache.Remember[[]Product](
        s.cache,
        key,
        func() ([]Product, error) {
            var ps []Product
            offset := (page - 1) * limit
            return ps, s.db.Offset(offset).Limit(limit).Find(&ps).Error
        },
        5*time.Minute,
    )
}
```

## What this module does NOT provide

- No `Tags(...).Set/Flush` API — use `Group` for namespacing.
- No `RememberForever` — pass a long explicit TTL or rely on no expiration.
- No `Increment`/`Decrement` — see the [KV module](../database/kv.md) for
  atomic counters.
- No `Lock`/`Unlock` (distributed locks) — implement via `kv.SetNX` if
  needed.
- The type is `*cache.Cache`, not `*cache.Client`.

## Best Practices

1. Set TTLs that match data volatility.
2. Use a `Group` per logical cache; flush via that module's cache instance.
3. Invalidate on writes to prevent stale data.
4. Handle cache misses gracefully — they are normal.
5. Don't cache sensitive data without a deliberate decision.

## Cache Warming

Pre-populate the cache on startup using the `Boot` hook:

```go
type CacheWarmer struct {
    cache *cache.Cache   `inject:""`
    svc   *ProductService `inject:""`
}

func (w *CacheWarmer) Boot(k types.Kernel) error {
    products, _ := w.svc.GetPopularProducts()
    for _, p := range products {
        _ = w.cache.Set("product:"+p.ID, p, time.Hour)
    }
    return nil
}
```

## Next Steps

- [Queues](queues.md) - Background job processing
- [KV Store](../database/kv.md) - Key-value operations and counters
