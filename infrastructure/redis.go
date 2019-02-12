package infrastructure

import (
	"github.com/go-redis/redis"
)

// NewRedisClient returns a client for redis.
func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client
}
