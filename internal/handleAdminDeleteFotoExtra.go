/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"errors"
	"net/http"
	"os"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/model"
)

func (s *Server) handleAdminDeleteFotoExtra() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(model.Session)

		uploadPath := os.Getenv("UPLOAD_PATH")

		fotoId := r.URL.Query().Get("foto")
		if fotoId == "" {
			logger.Error("falta el parámetro foto en la URL")
			s.handleError(w, 500, messages.ErrorFormulario)
			return
		}

		/*
			Comprobar que el usuario que intenta eliminar la fotografía sea el autor, y que la publicación no esté pública
		*/
		row, err := database.GetDB().Query(`SELECT COALESCE(fecha_publicacion, ""), autor_id FROM publicaciones WHERE id = ?`)
		if err != nil {
			logger.Error("error recuperando publicación cuya foto se va a borrar: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		var dbFechaPub, dbAutorId string
		row.Scan(&dbFechaPub, &dbAutorId)

		if dbFechaPub != "" {
			logger.Error("la publicación ya es pública, impidiendo borrado")
			s.handleError(w, 400, "No se puede eliminar fotografías de una publicación pública")
			return
		}

		// TODO: Crear permiso para saltarse esta limitación
		if sess.Autor_id != dbAutorId {
			logger.Error("el usuario no es el autor, impidiendo borrado")
			s.handleError(w, 400, messages.ErrorSinPermiso)
			return
		}

		/*
			Comprobar que la fotografía existe, y si es así eliminarla
		*/
		if f, err := os.Stat(uploadPath + "/extra/" + fotoId); err == nil {
			e2 := os.Remove(uploadPath + "/extra/" + f.Name())
			if e2 != nil {
				logger.Error("error borrando fotografía: %s", err.Error())
				s.handleError(w, 500, "Error borrando fotografía")
				return
			}
			w.Write([]byte("{ \"error\": false }"))
		} else if errors.Is(err, os.ErrNotExist) {
			logger.Error("la fotografía no existe: %s", err.Error())
			s.handleError(w, 404, messages.ErrorNoResultados)
		} else {
			logger.Error("error encontrando archivo a borrar: %s", err.Error())
			s.handleError(w, 500, messages.ErrorNoResultados)
		}
	}
}
