// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

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
	s.Router = s.SetupApiRoutes(s.Router)
	s.Router = s.JsonifyRoutes(router, "/api/v1")
	s.Router = s.IdentifyRequests(router)
	return s
}

type ridContextKey string
