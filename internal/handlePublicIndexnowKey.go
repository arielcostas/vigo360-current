package internal

import (
	"fmt"
	"net/http"
)

func (s *Server) handlePublicIndexnowKey() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "%s\n", r.URL.Path[1:len(r.URL.Path)-4])
	}
}
