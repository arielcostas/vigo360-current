/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/model"
	"vigo360.es/new/internal/templates"
)

// TODO: Get rid of this
type ResumenPost struct {
	Id                string
	Titulo            string
	Fecha_publicacion sql.NullString
	CantTags          int
	Publicado         bool
	Autor_id          string
	Autor_nombre      string
}

func (s *Server) handleAdminListPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(model.Session)
		db := database.GetDB()
		posts := []ResumenPost{}

		err := db.Select(&posts, `SELECT publicaciones.id, titulo, (fecha_publicacion < NOW() && fecha_publicacion IS NOT NULL) as publicado, DATE_FORMAT(fecha_publicacion,'%d-%m-%Y') as fecha_publicacion, autor_id, autores.nombre as autor_nombre, count(tag_id) as canttags FROM publicaciones LEFT JOIN autores ON publicaciones.autor_id = autores.id LEFT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id GROUP BY publicaciones.id ORDER BY publicado ASC, publicaciones.fecha_publicacion DESC;`)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando listado de publicacionoes: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		err = templates.Render(w, "admin-post.html", struct {
			Posts   []ResumenPost
			Session model.Session
		}{
			Posts:   posts,
			Session: sess,
		})

		if err != nil {
			logger.Error("error recuperando el autor: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
