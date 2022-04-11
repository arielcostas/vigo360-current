/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"net/http"

	"vigo360.es/new/internal/logger"
)

type appError struct {
	// Error whose message is extracted
	Error error
	// Message logged
	Message string
	// Response to give to the user
	Response string
	// HTTP Reply status code
	Status int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		logger.Error("[%s] %s: %s", r.URL.Path, err.Message, err.Error.Error())
		w.WriteHeader(err.Status)
		w.Write([]byte(err.Response))
		w.Write([]byte("\nSi crees que se trata de un error, contacta con el administrador."))
	}
}
