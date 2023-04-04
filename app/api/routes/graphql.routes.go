package routes

import (
	"github.com/go-chi/chi"
)

// GraphQL Routes function
func (rt *Routes) GraphQL(r chi.Router) {
	h := rt.Handlers.GraphQLHandler

	r.Route("/gql", func(r chi.Router) {
		r.Handle("/", h.Playground())
		r.Handle("/query", h.Query())
	})
}
