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

func (s *Server) handleAdminListSeries() http.HandlerFunc {
	type response struct {
		Series  []models.Serie
		Session models.Session
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)

		var series []models.Serie
		series, err := s.store.serie.Listar()
		// err := database.GetDB().Select(&series, `SELECT series.*, COUNT(publicaciones.id) as articulos FROM series LEFT JOIN publicaciones ON series.id = publicaciones.serie_id GROUP BY series.id;`)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando series: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		err = templates.Render(w, "admin-series.html", response{
			Series:  series,
			Session: sess,
		})
		if err != nil {
			logger.Error("error generando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
