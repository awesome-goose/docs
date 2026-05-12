# Rate Limiting

Protect your API from abuse with rate limiting.

## Overview

Rate limiting restricts the number of requests a client can make within a time window. This prevents:

- DoS attacks
- Brute force attempts
- API abuse
- Resource exhaustion

## Basic Rate Limiter

### Memory-Based Limiter

```go
type RateLimiter struct {
    requests map[string]*RateLimit
    mutex    sync.RWMutex
    limit    int
    window   time.Duration
}

type RateLimit struct {
    Count     int
    ExpiresAt time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
    return &RateLimiter{
        requests: make(map[string]*RateLimit),
        limit:    limit,
        window:   window,
    }
}

func (r *RateLimiter) IsAllowed(key string) bool {
    r.mutex.Lock()
    defer r.mutex.Unlock()

    now := time.Now()

    // Get or create rate limit
    rl, exists := r.requests[key]
    if !exists || now.After(rl.ExpiresAt) {
        r.requests[key] = &RateLimit{
            Count:     1,
            ExpiresAt: now.Add(r.window),
        }
        return true
    }

    // Check limit
    if rl.Count >= r.limit {
        return false
    }

    rl.Count++
    return true
}

func (r *RateLimiter) Remaining(key string) int {
    r.mutex.RLock()
    defer r.mutex.RUnlock()

    rl, exists := r.requests[key]
    if !exists {
        return r.limit
    }

    if time.Now().After(rl.ExpiresAt) {
        return r.limit
    }

    return r.limit - rl.Count
}
```

## SQL-Backed Rate Limiter

For distributed systems, store the counters in the `kv` module so all
instances share state via the SQL backend:

```go
import "github.com/awesome-goose/goose/modules/kv"

type DistributedRateLimiter struct {
    kv     *kv.KV `inject:""`
    limit  int
    window time.Duration
}

func NewDistributedRateLimiter(limit int, window time.Duration) *DistributedRateLimiter {
    return &DistributedRateLimiter{limit: limit, window: window}
}

func (r *DistributedRateLimiter) IsAllowed(key string) (bool, int, time.Time) {
    rateKey := "ratelimit:" + key

    // Increment counter (atomic in the kv module)
    count, _ := r.kv.Incr(rateKey)

    // Set expiry on first request
    if count == 1 {
        _, _ = r.kv.Expire(rateKey, r.window)
    }

    // Get TTL for reset time
    ttl, _ := r.kv.TTL(rateKey)
    resetAt := time.Now().Add(time.Duration(ttl) * time.Second)

    remaining := r.limit - int(count)
    if remaining < 0 {
        remaining = 0
    }

    return count <= int64(r.limit), remaining, resetAt
}
```

## Rate Limit Middleware

### Basic Middleware

```go
type RateLimitMiddleware struct {
    limiter *RateLimiter
}

func NewRateLimitMiddleware(limit int, window time.Duration) *RateLimitMiddleware {
    return &RateLimitMiddleware{
        limiter: NewRateLimiter(limit, window),
    }
}

func (m *RateLimitMiddleware) Handle(ctx types.Context, next types.Next) any {
    // Get client identifier
    key := m.getClientKey(ctx)

    // Check rate limit
    if !m.limiter.IsAllowed(key) {
        return ctx.Status(429).JSON(map[string]string{
            "error": "Too many requests",
        })
    }

    return next()
}

func (m *RateLimitMiddleware) getClientKey(ctx types.Context) string {
    // Try authenticated user first
    if userID := ctx.Get("user_id"); userID != nil {
        return "user:" + userID.(string)
    }

    // Fall back to IP address
    return "ip:" + ctx.IP()
}
```

### With Headers

```go
func (m *RateLimitMiddleware) Handle(ctx types.Context, next types.Next) any {
    key := m.getClientKey(ctx)

    allowed, remaining, resetAt := m.limiter.Check(key)

    // Set rate limit headers
    ctx.SetHeader("X-RateLimit-Limit", strconv.Itoa(m.limiter.limit))
    ctx.SetHeader("X-RateLimit-Remaining", strconv.Itoa(remaining))
    ctx.SetHeader("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

    if !allowed {
        retryAfter := int(time.Until(resetAt).Seconds())
        ctx.SetHeader("Retry-After", strconv.Itoa(retryAfter))

        return ctx.Status(429).JSON(map[string]interface{}{
            "error":       "Too many requests",
            "retry_after": retryAfter,
        })
    }

    return next()
}
```

## Rate Limit Strategies

### By IP Address

```go
func ByIP() KeyExtractor {
    return func(ctx types.Context) string {
        return "ip:" + ctx.IP()
    }
}
```

### By User

```go
func ByUser() KeyExtractor {
    return func(ctx types.Context) string {
        if userID := ctx.Get("user_id"); userID != nil {
            return "user:" + userID.(string)
        }
        return "ip:" + ctx.IP()
    }
}
```

### By API Key

```go
func ByAPIKey() KeyExtractor {
    return func(ctx types.Context) string {
        apiKey := ctx.Header("X-API-Key")
        if apiKey != "" {
            return "apikey:" + apiKey
        }
        return "ip:" + ctx.IP()
    }
}
```

### By Endpoint

```go
func ByEndpoint() KeyExtractor {
    return func(ctx types.Context) string {
        return "endpoint:" + ctx.Method() + ":" + ctx.Path() + ":ip:" + ctx.IP()
    }
}
```

## Configurable Rate Limiter

```go
type RateLimitConfig struct {
    Limit       int
    Window      time.Duration
    KeyFunc     func(ctx types.Context) string
    SkipFunc    func(ctx types.Context) bool
    ErrorFunc   func(ctx types.Context, limit int, remaining int, resetAt time.Time) any
}

func DefaultRateLimitConfig() *RateLimitConfig {
    return &RateLimitConfig{
        Limit:  100,
        Window: time.Minute,
        KeyFunc: func(ctx types.Context) string {
            return ctx.IP()
        },
        SkipFunc: nil,
        ErrorFunc: func(ctx types.Context, limit int, remaining int, resetAt time.Time) any {
            return ctx.Status(429).JSON(map[string]string{
                "error": "Too many requests",
            })
        },
    }
}

type ConfigurableRateLimiter struct {
    config  *RateLimitConfig
    limiter *RateLimiter
}

func (m *ConfigurableRateLimiter) Handle(ctx types.Context, next types.Next) any {
    // Check skip function
    if m.config.SkipFunc != nil && m.config.SkipFunc(ctx) {
        return next()
    }

    key := m.config.KeyFunc(ctx)
    allowed, remaining, resetAt := m.limiter.Check(key)

    // Set headers
    ctx.SetHeader("X-RateLimit-Limit", strconv.Itoa(m.config.Limit))
    ctx.SetHeader("X-RateLimit-Remaining", strconv.Itoa(remaining))
    ctx.SetHeader("X-RateLimit-Reset", strconv.FormatInt(resetAt.Unix(), 10))

    if !allowed {
        return m.config.ErrorFunc(ctx, m.config.Limit, remaining, resetAt)
    }

    return next()
}
```

## Different Limits for Different Routes

```go
func (c *Controller) Routes() types.Routes {
    // Strict limit for login (prevent brute force)
    loginLimit := NewRateLimitMiddleware(5, time.Minute)

    // Normal limit for API
    apiLimit := NewRateLimitMiddleware(100, time.Minute)

    // Higher limit for authenticated users
    authLimit := NewRateLimitMiddleware(1000, time.Minute)

    return types.Routes{
        {Method: "POST", Path: "/auth/login", Handler: c.Login,
            Middlewares: []any{loginLimit}},
        {Method: "GET", Path: "/api/public", Handler: c.PublicData,
            Middlewares: []any{apiLimit}},
        {Method: "GET", Path: "/api/users", Handler: c.Users,
            Middlewares: []any{&AuthMiddleware{}, authLimit}},
    }
}
```

## Tiered Rate Limits

Different limits for different user tiers:

```go
type TieredRateLimiter struct {
    limiters map[string]*RateLimiter
}

func NewTieredRateLimiter() *TieredRateLimiter {
    return &TieredRateLimiter{
        limiters: map[string]*RateLimiter{
            "free":       NewRateLimiter(100, time.Hour),
            "basic":      NewRateLimiter(1000, time.Hour),
            "pro":        NewRateLimiter(10000, time.Hour),
            "enterprise": NewRateLimiter(100000, time.Hour),
        },
    }
}

func (t *TieredRateLimiter) Handle(ctx types.Context, next types.Next) any {
    // Get user tier
    tier := ctx.Get("user_tier")
    if tier == nil {
        tier = "free"
    }

    limiter, ok := t.limiters[tier.(string)]
    if !ok {
        limiter = t.limiters["free"]
    }

    key := ctx.Get("user_id").(string)
    if !limiter.IsAllowed(key) {
        return ctx.Status(429).JSON(map[string]string{
            "error": "Rate limit exceeded for your plan",
        })
    }

    return next()
}
```

## Sliding Window via Per-Bucket Counters

The goose `kv` module is SQL-backed and does not expose sorted-set primitives
(`ZAdd`/`ZCard`/`ZRemRangeByScore`), so a true sorted-set sliding window is
not available. A reasonable approximation is to bucket counts by sub-window
and sum the recent buckets:

```go
type SlidingWindowLimiter struct {
    kv      *kv.KV `inject:""`
    limit   int
    window  time.Duration
    buckets int // e.g. 10 buckets per window for a 10x finer slide
}

func (l *SlidingWindowLimiter) IsAllowed(key string) bool {
    bucketDur := l.window / time.Duration(l.buckets)
    nowBucket := time.Now().Unix() / int64(bucketDur.Seconds())

    var total int64
    for i := 0; i < l.buckets; i++ {
        b := nowBucket - int64(i)
        rateKey := fmt.Sprintf("ratelimit:%s:%d", key, b)
        if v, err := l.kv.Get(rateKey); err == nil {
            if n, ok := v.(int64); ok {
                total += n
            }
        }
    }
    if total >= int64(l.limit) {
        return false
    }

    currentKey := fmt.Sprintf("ratelimit:%s:%d", key, nowBucket)
    if _, err := l.kv.Incr(currentKey); err == nil {
        _, _ = l.kv.Expire(currentKey, l.window+bucketDur)
    }
    return true
}
```

## Best Practices

1. **Use a distributed kv store** (Postgres/MySQL via the kv module) for
   multi-instance rate limiting.
2. **Set appropriate limits** based on your capacity
3. **Include rate limit headers** in responses
4. **Use different limits** for different endpoints
5. **Authenticate before rate limiting** where possible
6. **Log rate limit violations** for monitoring
7. **Provide clear error messages** with retry information
8. **Consider retry-after header** for better client experience

## Common Configurations

| Use Case        | Limit | Window     |
| --------------- | ----- | ---------- |
| Login attempts  | 5     | 15 minutes |
| Password reset  | 3     | 1 hour     |
| API (free tier) | 100   | 1 hour     |
| API (paid tier) | 1000  | 1 hour     |
| File uploads    | 10    | 1 hour     |
| Search queries  | 30    | 1 minute   |

## Next Steps

- [Authentication](authentication.md) - Secure your API
- [CORS](cors.md) - Cross-origin configuration
- [Security Overview](overview.md) - Security best practices
