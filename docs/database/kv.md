# Key-Value Store

The KV module provides Redis-compatible key-value storage for caching, sessions, and temporary data.

## Quick Start

```go
import "github.com/awesome-goose/goose/modules/kv"

// Configure the module
kvModule := kv.NewModule(
    kv.WithDriver("redis"),
    kv.WithHost("localhost"),
    kv.WithPort(6379),
)

// Include in application
stop, err := goose.Start(goose.API(platform, module, []types.Module{
    kvModule,
}))
```

## Configuration Options

### Redis

```go
kvModule := kv.NewModule(
    kv.WithDriver("redis"),
    kv.WithHost("localhost"),
    kv.WithPort(6379),
    kv.WithPassword(""),
    kv.WithDB(0),
)
```

### In-Memory (Development)

```go
kvModule := kv.NewModule(
    kv.WithDriver("memory"),
)
```

### Environment Configuration

```go
kvModule := kv.NewModule(
    kv.WithDriver(env.String("KV_DRIVER", "redis")),
    kv.WithHost(env.String("REDIS_HOST", "localhost")),
    kv.WithPort(env.Int("REDIS_PORT", 6379)),
    kv.WithPassword(env.String("REDIS_PASSWORD", "")),
)
```

**.env file:**

```env
KV_DRIVER=redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
```

## Injecting KV Client

Services receive the KV client via dependency injection:

```go
import "github.com/awesome-goose/goose/modules/kv"

type CacheService struct {
    kv *kv.Client `inject:""`
}
```

## Basic Operations

### Set Value

```go
func (s *CacheService) Set(key string, value interface{}) error {
    return s.kv.Set(key, value)
}
```

### Set with Expiration (TTL)

```go
func (s *CacheService) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
    return s.kv.SetEx(key, value, ttl)
}

// Example: Cache for 1 hour
s.kv.SetEx("user:123", user, time.Hour)

// Cache for 5 minutes
s.kv.SetEx("search:results:query", results, 5*time.Minute)
```

### Get Value

```go
func (s *CacheService) Get(key string) (string, error) {
    return s.kv.Get(key)
}

// Get with type conversion
func (s *CacheService) GetUser(id string) (*User, error) {
    data, err := s.kv.Get("user:" + id)
    if err != nil {
        return nil, err
    }

    var user User
    err = json.Unmarshal([]byte(data), &user)
    return &user, err
}
```

### Delete Value

```go
func (s *CacheService) Delete(key string) error {
    return s.kv.Del(key)
}

// Delete multiple keys
func (s *CacheService) DeleteMany(keys ...string) error {
    return s.kv.Del(keys...)
}
```

### Check Existence

```go
func (s *CacheService) Exists(key string) bool {
    exists, _ := s.kv.Exists(key)
    return exists
}
```

## Caching Patterns

### Cache-Aside

```go
func (s *ProductService) GetByID(id string) (*Product, error) {
    cacheKey := "product:" + id

    // Try cache first
    cached, err := s.kv.Get(cacheKey)
    if err == nil {
        var product Product
        json.Unmarshal([]byte(cached), &product)
        return &product, nil
    }

    // Fetch from database
    var product Product
    if err := s.db.First(&product, "id = ?", id).Error; err != nil {
        return nil, err
    }

    // Store in cache
    data, _ := json.Marshal(product)
    s.kv.SetEx(cacheKey, string(data), time.Hour)

    return &product, nil
}
```

### Write-Through

```go
func (s *ProductService) Update(id string, dto UpdateProductDTO) (*Product, error) {
    // Update database
    product, err := s.updateInDB(id, dto)
    if err != nil {
        return nil, err
    }

    // Update cache
    cacheKey := "product:" + id
    data, _ := json.Marshal(product)
    s.kv.SetEx(cacheKey, string(data), time.Hour)

    return product, nil
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
    s.kv.Del("product:" + id)

    // Invalidate related caches
    s.kv.Del("products:list")

    return nil
}
```

## Session Management

Store user sessions:

```go
type SessionService struct {
    kv *kv.Client `inject:""`
}

func (s *SessionService) Create(userID string) (string, error) {
    sessionID := uuid.New().String()

    session := map[string]interface{}{
        "user_id":    userID,
        "created_at": time.Now(),
    }

    data, _ := json.Marshal(session)
    err := s.kv.SetEx("session:"+sessionID, string(data), 24*time.Hour)

    return sessionID, err
}

func (s *SessionService) Get(sessionID string) (*Session, error) {
    data, err := s.kv.Get("session:" + sessionID)
    if err != nil {
        return nil, err
    }

    var session Session
    json.Unmarshal([]byte(data), &session)
    return &session, nil
}

func (s *SessionService) Destroy(sessionID string) error {
    return s.kv.Del("session:" + sessionID)
}

func (s *SessionService) Refresh(sessionID string) error {
    return s.kv.Expire("session:"+sessionID, 24*time.Hour)
}
```

## Rate Limiting

Implement rate limiting with KV:

```go
type RateLimiter struct {
    kv *kv.Client `inject:""`
}

func (r *RateLimiter) IsAllowed(key string, limit int, window time.Duration) bool {
    countKey := "ratelimit:" + key

    // Get current count
    countStr, err := r.kv.Get(countKey)
    if err != nil {
        // First request
        r.kv.SetEx(countKey, "1", window)
        return true
    }

    count, _ := strconv.Atoi(countStr)
    if count >= limit {
        return false
    }

    // Increment
    r.kv.Incr(countKey)
    return true
}
```

## Counter Operations

### Increment/Decrement

```go
func (s *StatsService) IncrementViews(postID string) (int64, error) {
    return s.kv.Incr("views:" + postID)
}

func (s *StatsService) DecrementStock(productID string) (int64, error) {
    return s.kv.Decr("stock:" + productID)
}

func (s *StatsService) IncrementBy(key string, amount int64) (int64, error) {
    return s.kv.IncrBy(key, amount)
}
```

## Hash Operations

Store structured data:

```go
// Set hash field
func (s *UserService) SetProfile(userID string, field string, value string) error {
    return s.kv.HSet("user:"+userID+":profile", field, value)
}

// Get hash field
func (s *UserService) GetProfileField(userID string, field string) (string, error) {
    return s.kv.HGet("user:"+userID+":profile", field)
}

// Get all hash fields
func (s *UserService) GetFullProfile(userID string) (map[string]string, error) {
    return s.kv.HGetAll("user:" + userID + ":profile")
}
```

## List Operations

Implement queues:

```go
// Push to list
func (s *QueueService) Push(queue string, value string) error {
    return s.kv.LPush(queue, value)
}

// Pop from list
func (s *QueueService) Pop(queue string) (string, error) {
    return s.kv.RPop(queue)
}

// Get list length
func (s *QueueService) Length(queue string) (int64, error) {
    return s.kv.LLen(queue)
}
```

## Set Operations

Track unique items:

```go
// Add to set
func (s *TagService) AddTag(postID string, tag string) error {
    return s.kv.SAdd("post:"+postID+":tags", tag)
}

// Get all set members
func (s *TagService) GetTags(postID string) ([]string, error) {
    return s.kv.SMembers("post:" + postID + ":tags")
}

// Check membership
func (s *TagService) HasTag(postID string, tag string) (bool, error) {
    return s.kv.SIsMember("post:"+postID+":tags", tag)
}
```

## Key Management

### Key Patterns

Use consistent key naming:

```go
// Pattern: resource:id:attribute
"user:123"                  // User data
"user:123:profile"          // User profile hash
"user:123:sessions"         // User sessions set
"product:456"               // Product data
"cache:products:list"       // Cached product list
"session:abc123"            // Session data
"ratelimit:ip:192.168.1.1"  // Rate limit counter
```

### Key Expiration

```go
// Set expiration
s.kv.Expire("temp:data", time.Hour)

// Get TTL
ttl, _ := s.kv.TTL("temp:data")

// Remove expiration
s.kv.Persist("temp:data")
```

### Pattern Matching

```go
// Get keys matching pattern
keys, _ := s.kv.Keys("user:*")

// Delete by pattern (use with caution)
for _, key := range keys {
    s.kv.Del(key)
}
```

## Best Practices

1. **Use namespaced keys** with consistent patterns
2. **Set TTL** for cached data to prevent stale reads
3. **Handle missing keys** gracefully
4. **Use appropriate data structures** (strings, hashes, lists, sets)
5. **Avoid storing large values** - prefer references
6. **Monitor memory usage** in production

## Error Handling

```go
func (s *CacheService) GetSafe(key string) (string, bool) {
    value, err := s.kv.Get(key)
    if err != nil {
        // Key doesn't exist or error
        return "", false
    }
    return value, true
}
```

## Next Steps

- [SQL Database](sql.md) - Relational database
- [Caching](../advanced/caching.md) - Advanced caching
- [Sessions](../building-blocks/sessions.md) - Session management
