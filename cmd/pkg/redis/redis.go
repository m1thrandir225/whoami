package redis

import "github.com/redis/go-redis/v9"

func NewRedisClient(redisURL string) (*redis.Client, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	return redis.NewClient(opt), nil
}
