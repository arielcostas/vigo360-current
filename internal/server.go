package internal

import (
	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	store  *Container
}

func NewServer(c *Container) *Server {
	s := &Server{
		store: c,
	}

	var router = mux.NewRouter().StrictSlash(true)
	router = s.SetupWebRoutes(router)
	router = s.JsonifyRoutes(router, "/api/v1")
	router = s.JsonifyRoutes(router, "/admin/async")

	router = s.IdentifyRequests(router)
	router = s.IdentifySessions(router)
	router = s.LogRequests(router)
	router = s.SetupSecurityHeaders(router)
	s.Router = router
	return s
}

type ridContextKey string
