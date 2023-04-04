package server

import (
	"gogql/app/api/dataloaders"
	"gogql/app/middlewares"
	"gogql/config"
	"time"

	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

type RestServer struct {
	Router *chi.Mux
}

// RestServer is ...
func NewRestServer(c *config.Clients) *RestServer {
	r := chi.NewRouter()
	restServer := &RestServer{r}

	// Add CORS
	corsOrigin(r)

	// read middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(middlewares.AuthTokenReader())
	r.Use(middlewares.OrgUIDReader())

	r.Route("/", func(r chi.Router) {
		urls(r, c)
	})

	return restServer
}

func (restServer *RestServer) Start(address string) {
	log.Printf("Listing and serving on port %s", address)
	http.ListenAndServe(address, restServer.Router)
}

func urls(r chi.Router, c *config.Clients) {
	dbStore, rt := Injection(c)

	r.Use(dataloaders.DataloaderMiddleware(dbStore))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	r.Route("/api", func(r chi.Router) {
		rt.Ping(r)
		rt.AuthRoutes(r)
		rt.GraphQL(r)
	})
}

// corsOrigin function
func corsOrigin(r *chi.Mux) {
	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:  []string{"*"},
		AllowOriginFunc: func(r *http.Request, origin string) bool { return true },
		AllowedMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Organization",
			"Owner",
			"Custodian",
			"Warehouse",
			"Proximity",
			"Content-Type",
			"X-CSRF-Token",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
}
