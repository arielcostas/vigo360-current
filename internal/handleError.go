// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"fmt"
	"net/http"
)

func (s *Server) handleError(w http.ResponseWriter, status int, message string) {
	// TODO: Implementar ID de petición
	// TODO: Mejorar esta página
	w.WriteHeader(status)
	fmt.Fprintf(w, `<!DOCTYPE html><html><body><h1>%s</h1></body></html>`, message)
}

func (s *Server) handleJsonError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, `{ "error": true, "errorMsg": "%s" }`, message)
}
