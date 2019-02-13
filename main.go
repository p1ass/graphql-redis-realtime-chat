package main

import (
	"github.com/naoki-kishi/graphql-redis-realtime-chat/infrastructure"
	"log"
)

func main() {

	redisURL := "redis:6379"

	client, err := infrastructure.NewRedisClient(redisURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	s := infrastructure.NewGraphQLServer(client)
	log.Fatal(s.Serve("/query", 8080))
}
