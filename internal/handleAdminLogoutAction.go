/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"net/http"
	"time"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
)

func revokeSession(sessid string) error {
	_, err := database.GetDB().Exec("UPDATE sesiones SET revocada = 1 WHERE sessid = ?;", sessid)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleAdminLogoutAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)
		if err := revokeSession(sess.Id); err != nil {
			logger.Error("error revocando sesión: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "sess",
			Value:    "",
			Path:     "/",
			Domain:   r.URL.Host,
			Expires:  time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC),
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
		})

		logger.Information("revoked session with id %s", sess.Id)
		w.Header().Add("Location", "/admin/login")
		w.WriteHeader(302)
	}
}
