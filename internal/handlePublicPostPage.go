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

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicPostPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		req_post_id := mux.Vars(r)["postid"]

		var post models.Publicacion
		if np, err := s.store.publicacion.ObtenerPorId(req_post_id, true); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró la publicación: %s", err.Error())
				s.handleError(w, 404, messages.ErrorPaginaNoEncontrada)
			} else {
				logger.Error("error recuperando la publicación: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}
			return
		} else {
			post = np
		}

		if post.Serie.Id != "" {
			var err error
			post.Serie, err = s.store.serie.Obtener(post.Serie.Id)
			if err != nil {
				logger.Error("error recuperando serie de la publicación: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}
		}

		var recommendations []Sugerencia
		if nr, err := generateSuggestions(post.Id, s.store.publicacion); err != nil {
			logger.Error("error recuperando sugerencias: %s", err.Error())
			recommendations = make([]Sugerencia, 0)
		} else {
			recommendations = nr
		}

		var keywords = ""
		for _, t := range post.Tags {
			keywords += t.Nombre + ","
		}

		post.Serie.Publicaciones = post.Serie.Publicaciones.FiltrarPublicas()

		if nc, err := s.store.comentario.ListarPublicos(post.Id); err != nil {
			logger.Error("error recuperando comentarios para %s: %s", post.Id, err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
		} else {
			post.Comentarios = nc
		}

		var err = templates.Render(w, "post-id.html", struct {
			Post            models.Publicacion
			Recommendations []Sugerencia
			Meta            PageMeta
		}{
			Post:            post,
			Recommendations: recommendations,
			Meta: PageMeta{
				Titulo:      post.Titulo,
				Descripcion: post.Resumen,
				Keywords:    keywords,
				Canonica:    fullCanonica("/post/" + post.Id),
				Miniatura:   fullCanonica("/static/thumb/" + post.Id + ".jpg"),
			},
		})
		if err != nil {
			logger.Error("error mostrando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
