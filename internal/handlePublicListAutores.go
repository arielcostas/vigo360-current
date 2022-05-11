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

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicListAutores() http.HandlerFunc {
	type Response struct {
		Autores []models.Autor
		Meta    PageMeta
	}

	var meta = PageMeta{
		Titulo:      "Autores",
		Descripcion: "Conoce a los autores y colaboradores de Vigo360.",
		Canonica:    fullCanonica("/autores"),
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		autores, err := s.store.autor.Listar()
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró ningún autor", err.Error())
				s.handleError(w, 404, "No se ha encontrado ningún autor ¿está el servidor bien configurado?")
			} else {
				logger.Error("error inesperado recuperando datos: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}
			return
		}

		err = templates.Render(w, "autores.html", Response{
			Autores: autores,
			Meta:    meta,
		})
		if err != nil {
			logger.Error("error renderizandoo plantilla: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
