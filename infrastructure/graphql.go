package infrastructure

import (
	"fmt"
	"github.com/99designs/gqlgen/handler"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/naoki-kishi/graphql-redis-realtime-chat/graphql"
	"github.com/rs/cors"
	"net/http"
	"sync"
)

type graphQLServer struct {
	redisClient     *redis.Client
	messageChannels map[string]chan graphql.Message
	mutex           sync.Mutex
}

// NewGraphQLServer returns GraphQL server.
func NewGraphQLServer(client *redis.Client) *graphQLServer {

	return &graphQLServer{}
}

func (s *graphQLServer) Serve(route string, port int) error {
	mux := http.NewServeMux()
	mux.Handle(
		route,
		handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: &graphql.Resolver{}}),
			handler.WebsocketUpgrader(websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			}),
		),
	)

	mux.Handle("/", handler.Playground("GraphQL playground", route))
	mux.Handle("/playground", handler.Playground("GraphQL playground", route))

	handler := cors.AllowAll().Handler(mux)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), handler)
}
