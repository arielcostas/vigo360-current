/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
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
