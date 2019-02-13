package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
	"log"
	"sync"
	"time"
)

// NewGraphQLConfig returns Config and start subscribing redis pubsub.
func NewGraphQLConfig(redisClient *redis.Client) Config {
	resolver := newResolver(redisClient)

	resolver.startSubscribingRedis()

	return Config{
		Resolvers: resolver,
	}
}

// Resolver implements ResolverRoot interface.
type Resolver struct {
	redisClient     *redis.Client
	messageChannels map[string]chan Message
	userChannels    map[string]chan string
	mutex           sync.Mutex
}

func newResolver(redisClient *redis.Client) *Resolver {
	return &Resolver{
		redisClient:     redisClient,
		messageChannels: map[string]chan Message{},
		userChannels:    map[string]chan string{},
		mutex:           sync.Mutex{},
	}
}

// Mutation returns a resolver for mutation.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns a resolver for query.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// Subscription returns a resolver for subsctiption.
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) PostMessage(ctx context.Context, user string, message string) (*Message, error) {
	isLogined, err := r.checkLogin(user)
	if err != nil {
		return nil, err
	}
	if !isLogined {
		return nil, errors.New("This user has not been created")
	}

	// extend session expire
	val, err := r.redisClient.SetXX(user, user, 60*time.Minute).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if val == false {
		return nil, errors.New("This user has not been created")
	}

	// Publish a message.
	m := Message{
		User:    user,
		Message: message,
	}
	mb, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
		return nil, err

	}
	r.redisClient.Publish("room", mb)

	log.Println("【Mutation】PostMessage : ", m)

	return &m, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, user string) (string, error) {
	// This means that users has to call CreateUser again in 60 minutes.
	val, err := r.redisClient.SetNX(user, user, 60*time.Minute).Result()
	if err != nil {
		log.Println(err)
		return "", err
	}
	if val == false {
		return "", errors.New("This User name has already used")
	}

	// Notify new user joined.
	r.mutex.Lock()
	for _, ch := range r.userChannels {
		ch <- user
	}
	r.mutex.Unlock()

	return user, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]string, error) {
	users, err := r.redisClient.Keys("*").Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	log.Println("【Query】Users : ", users)

	return users, nil

}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MessagePosted(ctx context.Context, user string) (<-chan Message, error) {
	isLogined, err := r.checkLogin(user)
	if err != nil {
		return nil, err
	}
	if !isLogined {
		return nil, errors.New("This user has not been created")
	}

	messageChan := make(chan Message, 1)
	r.mutex.Lock()
	r.messageChannels[user] = messageChan
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.messageChannels, user)
		r.mutex.Unlock()
		r.redisClient.Del(user)
	}()

	log.Println("【Subscription】MessagePosted : ", user)

	return messageChan, nil
}

func (r *subscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan string, error) {
	isLogined, err := r.checkLogin(user)
	if err != nil {
		return nil, err
	}
	if !isLogined {
		return nil, errors.New("This user has not been created")
	}

	userChan := make(chan string, 1)
	r.mutex.Lock()
	r.userChannels[user] = userChan
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.userChannels, user)
		r.mutex.Unlock()
		r.redisClient.Del(user)
	}()

	log.Println("【Subscription】UserJoined : ", user)

	return userChan, nil
}

func (r *Resolver) checkLogin(user string) (bool, error) {
	val, err := r.redisClient.Exists(user).Result()
	if err != nil {
		log.Println(err)
		return false, err
	}

	if val == 1 {
		return true, nil
	}
	return false, nil
}

func (r *Resolver) startSubscribingRedis() {
	log.Println("Start Subscribing Redis...")

	go func() {
		pubsub := r.redisClient.Subscribe("room")
		defer pubsub.Close()

		for {
			msgi, err := pubsub.Receive()
			if err != nil {
				panic(err)
			}

			switch msg := msgi.(type) {
			case *redis.Message:

				// Convert recieved string to Message.
				m := Message{}
				if err := json.Unmarshal([]byte(msg.Payload), &m); err != nil {
					log.Println(err)
					continue
				}

				// Notify new message.
				r.mutex.Lock()
				for _, ch := range r.messageChannels {
					ch <- m
				}
				r.mutex.Unlock()

			default:
			}
		}
	}()
}
