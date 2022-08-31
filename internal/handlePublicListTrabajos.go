// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

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

func (s *Server) handlePublicListTrabajos() http.HandlerFunc {
	type Response struct {
		Trabajos models.Trabajos
		Meta     PageMeta
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		trabajos, err := s.store.trabajo.Listar()

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error obteniendo trabajos: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		err = templates.Render(w, "trabajos.html", Response{
			Trabajos: trabajos,
			Meta: PageMeta{
				Titulo:      "Trabajos",
				Descripcion: "Trabajos originales e interesantes publicados por los autores de Vigo360.",
				Canonica:    fullCanonica("/trabajos"),
			},
		})

		if err != nil {
			logger.Error("error mostrando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
			return
		}

	}
}
