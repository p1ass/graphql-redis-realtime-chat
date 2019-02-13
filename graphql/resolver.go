package graphql

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"sync"
)

// Resolver implements ResolverRoot interface.
type Resolver struct {
	redisClient     *redis.Client
	messageChannels map[string]chan Message
	userChannels    map[string]chan string
	mutex           sync.Mutex
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

	err := r.createUser(user)
	if err != nil {
		return nil, err
	}

	// Publish a message.
	m := Message{
		User:    user,
		Message: message,
	}
	mb, _ := json.Marshal(m)
	r.redisClient.Publish("room", mb)

	return &m, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Users(ctx context.Context) ([]string, error) {
	cmd := r.redisClient.SMembers("users")
	if cmd.Err() != nil {
		log.Println(cmd.Err())
		return nil, cmd.Err()
	}

	res, err := cmd.Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return res, nil

}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MessagePosted(ctx context.Context, user string) (<-chan Message, error) {
	err := r.createUser(user)
	if err != nil {
		return nil, err
	}

	messages := make(chan Message, 1)
	r.mutex.Lock()
	r.messageChannels[user] = messages
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.messageChannels, user)
	}()

	return messages, nil
}
func (r *subscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan string, error) {
	err := r.createUser(user)
	if err != nil {
		return nil, err
	}

	users := make(chan string, 1)
	r.mutex.Lock()
	r.userChannels[user] = users
	r.mutex.Unlock()

	go func() {
		<-ctx.Done()
		r.mutex.Lock()
		delete(r.userChannels, user)
		r.mutex.Unlock()
	}()

	return users, nil
}

func (r *Resolver) createUser(user string) error {
	// Set user to redis list.
	if err := r.redisClient.SAdd("users", user).Err(); err != nil {
		log.Println(err)
		return err
	}

	// Notify new user joined.
	r.mutex.Lock()
	for _, ch := range r.userChannels {
		ch <- user
	}
	r.mutex.Unlock()
	return nil
}
