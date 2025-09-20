package security

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
}

type RateLimitConfig struct {
	Requests int           `json:"requests"`
	Window   time.Duration `json:"window"`
}

type RateLimitResult struct {
	Allowed    bool  `json:"allowed"`
	Remaining  int64 `json:"remaining"`
	ResetTime  int64 `json:"reset_time"`
	RetryAfter int64 `json:"retry_after,omitempty"`
}

var (
	DefaultRateLimit = RateLimitConfig{
		Requests: 1000,
		Window:   time.Hour * 1,
	}

	AuthRateLimit = RateLimitConfig{
		Requests: 50,
		Window:   15 * time.Minute,
	}

	RegistrationRateLimit = RateLimitConfig{
		Requests: 3,
		Window:   time.Hour * 1,
	}

	PasswordResetRateLimit = RateLimitConfig{
		Requests: 3,
		Window:   time.Hour * 1,
	}
)

func NewRateLimiter(redisClient *redis.Client) (*RateLimiter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &RateLimiter{
		redisClient: redisClient,
	}, nil
}

func (rl *RateLimiter) CheckRateLimit(ctx context.Context, key string, config RateLimitConfig) (*RateLimitResult, error) {
	now := time.Now()

	windowStart := now.Add(-config.Window)

	pipe := rl.redisClient.Pipeline()

	pipe.ZRemRangeByScore(ctx, key, "0", strconv.FormatInt(windowStart.Unix(), 10))

	countCmd := pipe.ZCard(ctx, key)

	pipe.ZAdd(ctx, key, redis.Z{
		Score:  float64(now.Unix()),
		Member: now.UnixNano(),
	})

	pipe.Expire(ctx, key, config.Window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to execute Redis pipeline: %v", err)
	}

	count := countCmd.Val()
	allowed := count < int64(config.Requests)

	result := &RateLimitResult{
		Allowed:   allowed,
		Remaining: int64(config.Requests) - count,
		ResetTime: now.Add(config.Window).Unix(),
	}

	if !allowed {
		// Calculate retry after time
		oldestRequest, err := rl.redisClient.ZRange(ctx, key, 0, 0).Result()
		if err == nil && len(oldestRequest) > 0 {
			if oldestTime, err := strconv.ParseInt(oldestRequest[0], 10, 64); err == nil {
				retryAfter := time.Unix(oldestTime, 0).Add(config.Window).Unix() - now.Unix()
				if retryAfter > 0 {
					result.RetryAfter = retryAfter
				}
			}
		}
	}

	return result, nil
}

func GetClientIP(ctx *gin.Context) string {
	if ip := ctx.GetHeader("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := ctx.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := ctx.GetHeader("CF-Connecting-IP"); ip != "" {
		return ip
	}

	return ctx.ClientIP()
}

func (rl *RateLimiter) RateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := GetClientIP(ctx)
		key := fmt.Sprintf("rate_limit:%s:%s", ctx.Request.URL.Path, ip)
		result, err := rl.CheckRateLimit(ctx, key, config)
		if err != nil {
			// If Redis is down, allow the request but log the error
			ctx.Next()
			return
		}

		ctx.Header("X-RateLimit-Limit", strconv.Itoa(config.Requests))
		ctx.Header("X-RateLimit-Remaining", strconv.FormatInt(result.Remaining, 10))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(result.ResetTime, 10))

		if !result.Allowed {
			ctx.Header("Retry-After", strconv.FormatInt(result.RetryAfter, 10))
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "Rate limit exceeded",
				"retry_after": result.RetryAfter,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (rl *RateLimiter) UserRateLimitMiddleware(config RateLimitConfig) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		payload, exists := ctx.Get("authorization_payload")
		if !exists {
			// If no user payload, fall back to IP-based rate limiting
			clientIP := GetClientIP(ctx)
			key := fmt.Sprintf("rate_limit:user:%s:%s", ctx.Request.URL.Path, clientIP)

			result, err := rl.CheckRateLimit(ctx, key, config)
			if err != nil {
				ctx.Next()
				return
			}

			if !result.Allowed {
				ctx.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
				})
				ctx.Abort()
				return
			}
		} else {
			// Use user ID for rate limiting
			userPayload, ok := payload.(*Payload)
			if !ok {
				ctx.Next()
				return
			}

			key := fmt.Sprintf("rate_limit:user:%d:%s", userPayload.UserID, ctx.Request.URL.Path)

			result, err := rl.CheckRateLimit(ctx, key, config)
			if err != nil {
				ctx.Next()
				return
			}

			if !result.Allowed {
				ctx.JSON(http.StatusTooManyRequests, gin.H{
					"error": "Rate limit exceeded",
				})
				ctx.Abort()
				return
			}
		}

		ctx.Next()
	}
}

func (rl *RateLimiter) ResetRateLimit(ctx context.Context, key string) error {
	return rl.redisClient.Del(ctx, key).Err()
}

func (rl *RateLimiter) ResetAllRateLimits(ctx context.Context) error {
	pattern := "rate_limit:*"
	keys, err := rl.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return rl.redisClient.Del(ctx, keys...).Err()
	}
	return nil
}

func (rl *RateLimiter) ResetRateLimitByIP(ctx context.Context, ip string) error {
	pattern := fmt.Sprintf("rate_limit:*:%s", ip)
	keys, err := rl.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return rl.redisClient.Del(ctx, keys...).Err()
	}
	return nil
}

func (rl *RateLimiter) ResetRateLimitByUser(ctx context.Context, userID int64) error {
	pattern := fmt.Sprintf("rate_limit:user:%d:*", userID)
	keys, err := rl.redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return rl.redisClient.Del(ctx, keys...).Err()
	}
	return nil
}

func (rl *RateLimiter) Close() error {
	return rl.redisClient.Close()
}

func (rl *RateLimiter) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return rl.redisClient.Ping(ctx).Err()
}
