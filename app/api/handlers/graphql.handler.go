package handlers

import (
	"gogql/app/api/graphql/generated/graph"
	"gogql/app/api/resolvers"
	"gogql/app/services"
	"gogql/app/store/filestore"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

type GraphQLHandler struct {
	services  *services.Services
	filestore *filestore.FileStore
}

var _ GraphQLHandlerInterface = &GraphQLHandler{}

type GraphQLHandlerInterface interface {
	Playground() http.HandlerFunc
	Query() *handler.Server
}

func NewGraphQLHandler(s *services.Services, fs *filestore.FileStore) *GraphQLHandler {
	return &GraphQLHandler{s, fs}
}

// Playground Handler
func (h *GraphQLHandler) Playground() http.HandlerFunc {
	return playground.Handler("GraphQL playground", "/api/gql/query")
}

// Query Handler
func (h *GraphQLHandler) Query() *handler.Server {
	resolvers := resolvers.NewResolver(h.services, h.filestore)
	return handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolvers}))
}
