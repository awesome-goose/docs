# Caching

Improve performance with the Goose caching module.

## Overview

The Cache module provides a SQL-backed caching interface with support for:

- TTL (Time-To-Live) support
- Cache groups for organization
- Automatic cleanup of expired entries

## Quick Start

```go
import "github.com/awesome-goose/goose/modules/cache"

// Configure cache module
cacheModule := cache.NewModule(
    cache.WithGroup("app"),
    cache.WithDefaultTTL(time.Hour),
    cache.WithCleanupInterval(10 * time.Minute),
)

// Include in application
stop, err := goose.Start(goose.API(platform, module, []types.Module{
    cacheModule,
}))
```

## Configuration

### Available Options

```go
cacheModule := cache.NewModule(
    // Group/namespace for cache entries
    cache.WithGroup(env.String("CACHE_GROUP", "default")),
    // Default TTL for entries (default: 5 minutes)
    cache.WithDefaultTTL(time.Hour),
    // How often to clean up expired entries (0 = no cleanup)
    cache.WithCleanupInterval(10 * time.Minute),
)
```

### Configuration Struct

```go
type Config struct {
    // Group is the namespace for cache entries
    Group string
    // DefaultTTL is the default time-to-live for cache entries
    DefaultTTL time.Duration
    // CleanupInterval is how often expired entries are cleaned up
    CleanupInterval time.Duration
}
```

## Using the Cache

### Inject Cache Client

```go
type ProductService struct {
    cache *cache.Client `inject:""`
    db    *gorm.DB      `inject:""`
}
```

### Basic Operations

```go
// Set value
func (s *ProductService) CacheProduct(product *Product) error {
    return s.cache.Set("product:"+product.ID, product, time.Hour)
}

// Get value
func (s *ProductService) GetCachedProduct(id string) (*Product, error) {
    var product Product
    err := s.cache.Get("product:"+id, &product)
    if err != nil {
        return nil, err
    }
    return &product, nil
}

// Delete value
func (s *ProductService) InvalidateProduct(id string) error {
    return s.cache.Delete("product:" + id)
}

// Check existence
func (s *ProductService) IsCached(id string) bool {
    return s.cache.Has("product:" + id)
}
```

## Caching Patterns

### Cache-Aside (Lazy Loading)

```go
func (s *ProductService) GetByID(id string) (*Product, error) {
    cacheKey := "product:" + id

    // Try cache first
    var product Product
    if err := s.cache.Get(cacheKey, &product); err == nil {
        return &product, nil
    }

    // Cache miss - fetch from database
    if err := s.db.First(&product, "id = ?", id).Error; err != nil {
        return nil, err
    }

    // Store in cache
    s.cache.Set(cacheKey, &product, time.Hour)

    return &product, nil
}
```

### Write-Through

```go
func (s *ProductService) Update(id string, dto UpdateProductDTO) (*Product, error) {
    // Update database
    var product Product
    if err := s.db.First(&product, "id = ?", id).Error; err != nil {
        return nil, err
    }

    product.Name = dto.Name
    product.Price = dto.Price

    if err := s.db.Save(&product).Error; err != nil {
        return nil, err
    }

    // Update cache immediately
    s.cache.Set("product:"+id, &product, time.Hour)

    return &product, nil
}
```

### Cache Invalidation

```go
func (s *ProductService) Delete(id string) error {
    // Delete from database
    if err := s.db.Delete(&Product{}, "id = ?", id).Error; err != nil {
        return err
    }

    // Invalidate cache
    s.cache.Delete("product:" + id)

    // Also invalidate list cache
    s.cache.Delete("products:list")

    return nil
}
```

## Remember Pattern

Fetch from cache or execute callback:

```go
func (s *ProductService) GetByID(id string) (*Product, error) {
    var product Product

    err := s.cache.Remember("product:"+id, time.Hour, &product, func() (interface{}, error) {
        var p Product
        err := s.db.First(&p, "id = ?", id).Error
        return &p, err
    })

    if err != nil {
        return nil, err
    }

    return &product, nil
}
```

### Forever Remember

Cache indefinitely (until manually invalidated):

```go
func (s *ConfigService) GetSettings() (*Settings, error) {
    var settings Settings

    err := s.cache.RememberForever("app:settings", &settings, func() (interface{}, error) {
        return s.loadSettingsFromDB()
    })

    return &settings, err
}
```

## Cache Tags

Group related cache entries for bulk invalidation:

```go
// Set with tags
func (s *ProductService) CacheProduct(product *Product) error {
    return s.cache.Tags("products", "category:"+product.CategoryID).Set(
        "product:"+product.ID,
        product,
        time.Hour,
    )
}

// Invalidate by tag
func (s *ProductService) InvalidateProductsByCategory(categoryID string) error {
    return s.cache.Tags("category:" + categoryID).Flush()
}

// Invalidate all products
func (s *ProductService) InvalidateAllProducts() error {
    return s.cache.Tags("products").Flush()
}
```

## Caching Lists

```go
func (s *ProductService) GetAll() ([]Product, error) {
    cacheKey := "products:list"

    // Try cache
    var products []Product
    if err := s.cache.Get(cacheKey, &products); err == nil {
        return products, nil
    }

    // Fetch from database
    s.db.Find(&products)

    // Cache the list
    s.cache.Set(cacheKey, products, 30*time.Minute)

    return products, nil
}
```

## Pagination Caching

```go
func (s *ProductService) GetPaginated(page, limit int) ([]Product, error) {
    cacheKey := fmt.Sprintf("products:page:%d:limit:%d", page, limit)

    var products []Product
    if err := s.cache.Get(cacheKey, &products); err == nil {
        return products, nil
    }

    // Fetch from database
    offset := (page - 1) * limit
    s.db.Offset(offset).Limit(limit).Find(&products)

    // Cache with shorter TTL for paginated results
    s.cache.Set(cacheKey, products, 5*time.Minute)

    return products, nil
}
```

## HTTP Response Caching

Cache entire API responses:

```go
type CacheMiddleware struct {
    cache *cache.Client `inject:""`
    ttl   time.Duration
}

func (m *CacheMiddleware) Handle(ctx types.Context, next types.Next) any {
    // Only cache GET requests
    if ctx.Method() != "GET" {
        return next()
    }

    cacheKey := "response:" + ctx.Path()

    // Try cache
    var cached map[string]interface{}
    if err := m.cache.Get(cacheKey, &cached); err == nil {
        ctx.SetHeader("X-Cache", "HIT")
        return cached
    }

    // Execute request
    response := next()

    // Cache response
    m.cache.Set(cacheKey, response, m.ttl)
    ctx.SetHeader("X-Cache", "MISS")

    return response
}
```

## Atomic Operations

```go
// Increment counter
func (s *StatsService) IncrementViews(postID string) (int64, error) {
    return s.cache.Increment("views:"+postID, 1)
}

// Decrement counter
func (s *StatsService) DecrementStock(productID string) (int64, error) {
    return s.cache.Decrement("stock:"+productID, 1)
}
```

## Distributed Locking

Prevent concurrent operations:

```go
func (s *OrderService) ProcessOrder(orderID string) error {
    lockKey := "lock:order:" + orderID

    // Acquire lock
    acquired, err := s.cache.Lock(lockKey, 30*time.Second)
    if err != nil || !acquired {
        return fmt.Errorf("order already being processed")
    }
    defer s.cache.Unlock(lockKey)

    // Process order safely
    return s.doProcessOrder(orderID)
}
```

## Best Practices

1. **Set appropriate TTLs** based on data volatility
2. **Use cache prefixes** to namespace keys
3. **Invalidate on writes** to prevent stale data
4. **Use tags** for related cache invalidation
5. **Handle cache misses gracefully**
6. **Monitor hit/miss ratios**
7. **Don't cache sensitive data** unnecessarily

## Cache Warming

Pre-populate cache on startup:

```go
type CacheWarmer struct {
    cache          *cache.Client   `inject:""`
    productService *ProductService `inject:""`
}

func (w *CacheWarmer) OnStart() {
    // Warm frequently accessed data
    products, _ := w.productService.GetPopularProducts()
    for _, product := range products {
        w.cache.Set("product:"+product.ID, product, time.Hour)
    }
}
```

## Next Steps

- [Queues](queues.md) - Background job processing
- [KV Store](../database/kv.md) - Key-value operations
- [Performance](../deployment/performance.md) - Optimization tips
