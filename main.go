package main

import (
	"bufio"
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"time"
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

func main() {

	if len(os.Args) != 2 {
		fmt.Println("There is an error in the argument")
		os.Exit(1)
	}
	userName := os.Args[1]
	userKey := "online" + userName

	client := NewRedisClient()
	defer client.Close()

	val, err := client.SetNX(userKey, userName, 2*time.Minute).Result()
	defer client.Del(userKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if val == false {
		fmt.Println("User already online")
		os.Exit(1)
	}

	// 受け取ったメッセージをsubChanに流すgoroutine
	subChan := make(chan string)

	go func() {
		clientForRoom := NewRedisClient()
		defer clientForRoom.Close()

		pubsub := clientForRoom.Subscribe("room")
		defer pubsub.Close()

		for {
			msgi, err := pubsub.Receive()
			if err != nil {
				panic(err)
			}

			switch msg := msgi.(type) {
			case *redis.Subscription:
				break
			case *redis.Message:
				subChan <- fmt.Sprint(msg.Payload)

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

	chatExit := false

	for !chatExit {
		prompt := ">"
		fmt.Print(prompt)

		select {
		case msg := <-subChan:
			fmt.Println(msg)

		case line := <-sayChan:
			if line == "/exit" {
				chatExit = true
				return
			}
			client.Publish("room", userName+": "+line)
		}
	}
}
