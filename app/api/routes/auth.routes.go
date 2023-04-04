package routes

import (
	"github.com/go-chi/chi"
)

// AuthRoutes Routes function
func (rt *Routes) AuthRoutes(r chi.Router) {
	h := rt.Handlers.AuthHandler

	r.Route("/auth", func(r chi.Router) {
		r.Get("/permissions", h.ListPermissions)
		r.Post("/logout", h.Logout)
	})
}
