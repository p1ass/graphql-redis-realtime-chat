package main

import (
	"bufio"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/naoki-kishi/graphql-redis-realtime-chat/infrastructure"
	"log"
	"os"
	"time"
)

func main() {

	redisURL := "redis:6379"

	// if len(os.Args) != 2 {
	// 	fmt.Println("There is an error in the argument")
	// 	os.Exit(1)
	// }
	userName := "fuga"
	userKey := "online" + userName

	client, err := infrastructure.NewRedisClient(redisURL)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	val, err := client.SetNX(userKey, userName, 2*time.Minute).Result()
	defer client.Del(userKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	val, err := client.SAdd("users", userName)
	if err:= nil{
		panic(err)
	}

	if val == false {
		fmt.Println("User already online")
		os.Exit(1)
	}

	// 受け取ったメッセージをmsgChanに流すgoroutine
	msgChan := make(chan string)

	go func() {
		clientForRoom, err := infrastructure.NewRedisClient(redisURL)
		if err != nil {
			panic(err)
		}
		defer clientForRoom.Close()

		pubsub := client.Subscribe("room")
		defer pubsub.Close()

		for {
			msgi, err := pubsub.Receive()
			if err != nil {
				panic(err)
			}

			switch msg := msgi.(type) {
			case *redis.Subscription:
				msgChan <- "connected"

			case *redis.Message:
				msgChan <- fmt.Sprint(msg.Payload)

			default:
				panic("unreached")
			}
		}
	}()

	// コマンドプロンプトで受け取った入力をsayChanに流すgoroutine
	sayChan := make(chan string)
	go func() {
		bio := bufio.NewReader(os.Stdin)

		for {
			line, _, err := bio.ReadLine()
			if err != nil {
				fmt.Println(err)
				sayChan <- "/exit"
				return
			}
			sayChan <- string(line)
		}
	}()

	client.Publish("room", userName+" has joined")
	defer client.Publish("room", userName+" has left")

	// 各channelから受け取った情報を出力する部分
	// chatExit := false
	// for !chatExit {
	// 	prompt := ">"
	// 	fmt.Print(prompt)

	// 	select {
	// 	case msg := <-msgChan:
	// 		fmt.Println(msg)

	// 	case line := <-sayChan:
	// 		if line == "/exit" {
	// 			chatExit = true
	// 			return
	// 		}
	// 		client.Publish("room", userName+": "+line)
	// 	}
	// }

	s := infrastructure.NewGraphQLServer(client)
	log.Fatal(s.Serve("/query", 8080))
}
