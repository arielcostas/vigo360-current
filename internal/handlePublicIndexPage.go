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
	"strconv"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

type indexParams struct {
	CurrentPage int
	PageCount   int
	Posts       models.Publicaciones
	Meta        PageMeta
}

func (s *Server) handlePublicIndex() http.HandlerFunc {
	var meta = PageMeta{
		Titulo:      "Inicio",
		Descripcion: "Vigo360 es un proyecto dedicado a estudiar varios aspectos de la ciudad de Vigo (España) y su área de influencia, centrándose en la toponimia y el transporte.",
		Canonica:    fullCanonica("/"),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		posts, err := s.store.publicacion.Listar()
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando datos: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		posts = posts.FiltrarPublicas()

		/* Paginación */
		var pagina = 1

		var queryPage = r.URL.Query().Get("page")
		if queryPage != "" {
			o, err := strconv.Atoi(queryPage)
			if err != nil {
				logger.Error("no se pudo convertir '%s' a un número de página", queryPage)
				s.handleError(w, 404, messages.ErrorNoResultados)
				return
			}
			pagina = o
		}

		var inicio = pagina*9 - 9
		var limite = getMinimo(inicio+9, len(posts))

		if inicio >= len(posts) || inicio < 0 {
			logger.Error("con %d publicaciones no existe la página %s", len(posts), pagina)
			s.handleError(w, 404, messages.ErrorNoResultados)
			return
		}

		var restantes = len(posts) - 9 // Los artículos que aún no se han metido en una página
		var cantidadPaginas = 1
		for restantes > 0 {
			cantidadPaginas++
			restantes -= 9
		}

		err = templates.Render(w, "index.html", indexParams{
			CurrentPage: pagina,
			PageCount:   cantidadPaginas,
			Posts:       posts[inicio:limite],
			Meta:        meta,
		})
		if err != nil {
			logger.Error("error renderizando la página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
