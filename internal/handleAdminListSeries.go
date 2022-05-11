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

func (s *Server) handleAdminListSeries() http.HandlerFunc {
	type response struct {
		Series  []Serie
		Session model.Session
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(model.Session)

		series := []Serie{}
		err := database.GetDB().Select(&series, `SELECT series.*, COUNT(publicaciones.id) as articulos FROM series LEFT JOIN publicaciones ON series.id = publicaciones.serie_id GROUP BY series.id;`)
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
