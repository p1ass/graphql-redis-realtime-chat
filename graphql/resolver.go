package graphql

import (
	"context"
)

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) Subscription() SubscriptionResolver {
	return &subscriptionResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) PostMessage(ctx context.Context, user string, message string) (*Message, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Messages(ctx context.Context) ([]Message, error) {
	panic("not implemented")
}
func (r *queryResolver) Users(ctx context.Context) ([]string, error) {
	panic("not implemented")
}

type subscriptionResolver struct{ *Resolver }

func (r *subscriptionResolver) MessagePosted(ctx context.Context, user string) (<-chan Message, error) {
	panic("not implemented")
}
func (r *subscriptionResolver) UserJoined(ctx context.Context, user string) (<-chan string, error) {
	panic("not implemented")
}
