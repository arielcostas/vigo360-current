/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"fmt"
	"net/http"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/service"
)

func (s *Server) handleAdminRechazarComentario() http.HandlerFunc {
	var cs = service.NewComentarioService(s.store.comentario, s.store.publicacion)

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)
		cid := r.URL.Query().Get("cid") // cid = commentId = el comentario a rechazar

		err := cs.Rechazar(cid, sess.Autor_id)
		if err != nil {
			logger.Error("error rechazando comentario %s: %s", cid, err.Error())
			fmt.Fprintf(w, "Hubo un error rechazando el comentario")
		}

		w.Header().Add("Location", "/admin/comentarios")
		defer w.WriteHeader(303)
	}
}
