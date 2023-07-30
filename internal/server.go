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
	s.Router = s.SetupWebRoutes(router)
	s.Router = s.JsonifyRoutes(router, "/api/v1")
	s.Router = s.JsonifyRoutes(router, "/admin/async")
	s.Router = s.IdentifyRequests(router)
	return s
}

type ridContextKey string
