package routes

import (
	"gogql/app/api/handlers"
)

type Routes struct {
	Handlers *handlers.Handlers
}

func NewRoutes(h *handlers.Handlers) *Routes {
	return &Routes{h}
}
