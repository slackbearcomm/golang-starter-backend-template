package routes

import (
	"gogql/app/api/handlers"

	"github.com/go-chi/chi"
)

// Ping Routes function
func (rt *Routes) Ping(r chi.Router) {
	r.Route("/", func(r chi.Router) {
		r.Get("/ping", handlers.Ping)
		r.Get("/king", handlers.King)
		r.Get("/ding", handlers.Ding)
	})
}
