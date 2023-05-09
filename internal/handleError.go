package internal

import (
	"fmt"
	"net/http"

	"vigo360.es/new/internal/messages"
)

func (s *Server) handleError(w http.ResponseWriter, status int, message messages.ErrorMessage) {
	// TODO: Implementar ID de petición
	// TODO: Mejorar esta página
	w.WriteHeader(status)
	fmt.Fprintf(w, `<!DOCTYPE html><html><body><h1>%s</h1></body></html>`, message)
}

func (s *Server) handleJsonError(w http.ResponseWriter, status int, message messages.ErrorMessage) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{ "error": "%s" }`, message)
}
