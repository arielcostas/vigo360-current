package internal

import (
	"fmt"
	"net/http"

	"vigo360.es/new/internal/messages"
)

func (s *Server) handleError(r *http.Request, w http.ResponseWriter, status int, message messages.ErrorMessage) {
	w.WriteHeader(status)
	rid := r.Context().Value(ridContextKey("rid")).(string)

	_, err := fmt.Fprintf(w, `<!DOCTYPE html><html><body><h1>%s</h1></body></html>`, message)
	if err != nil {
		_, _ = fmt.Fprintf(w, `{ "error": "%s" }`, messages.ErrorFatal)
		fmt.Println(rid, err)
	}
}

func (s *Server) handleJsonError(r *http.Request, w http.ResponseWriter, status int, message messages.ErrorMessage) {
	w.WriteHeader(status)
	rid := r.Context().Value(ridContextKey("rid")).(string)
	_, err := fmt.Fprintf(w, `{ "error": "%s", "rid": "%s" }`, message, rid)
	if err != nil {
		_, _ = fmt.Fprintf(w, `{ "error": "%s", "rid": "%s" }`, messages.ErrorFatal, rid)
		fmt.Println(rid, err)
	}
}
