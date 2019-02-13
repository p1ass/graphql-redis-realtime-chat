package infrastructure

import (
	"github.com/go-redis/redis"
)

// NewRedisClient returns a client for redis.
func NewRedisClient(url string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	return client, err
}
