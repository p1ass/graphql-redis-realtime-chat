package main

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/naoki-kishi/graphql-redis-realtime-chat/infrastructure"
	"log"
)

type config struct {
	RedisURL string `envconfig:"REDIS_URL"`
	Port     int    `envconfig:"PORT"`
}

func main() {

	var config config
	err := envconfig.Process("", &config)
	if err != nil {
		log.Println(err)
	}

	client, err := infrastructure.NewRedisClient(config.RedisURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	s := infrastructure.NewGraphQLServer(client)
	log.Fatal(s.Serve("/query", config.Port))
}
