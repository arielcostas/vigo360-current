/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"net/http"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/service"
)

func (s *Server) handlePublicEnviarComentario() http.HandlerFunc {
	var cs = service.NewComentarioService(s.store.comentario, s.store.publicacion)

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))

		r.ParseForm()
		var publicacion_id = mux.Vars(r)["postid"]
		var nombre = r.Form.Get("nombre")
		var contenido = r.Form.Get("contenido")

		nc, err := cs.AgregarComentario(publicacion_id, nombre, contenido)
		if err != nil {
			logger.Error("error guardando comentario: %s", err.Error())
			s.handleError(w, 400, err.Error())
			return
		}

		logger.Information("guardado comentario con ID %s", nc.Id)

		w.Header().Add("Location", r.URL.Path)
		defer w.WriteHeader(http.StatusSeeOther)
	}
}
