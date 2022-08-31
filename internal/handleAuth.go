// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"context"
	"net/http"
	"net/url"

	"vigo360.es/new/internal/logger"
)

type sessionContextKey string

func (s *Server) withAuth(h http.HandlerFunc) http.HandlerFunc {
	var gotoLogin = func(w http.ResponseWriter, rawnext string) {
		w.Header().Add("Location", "/admin/login?next="+url.QueryEscape(rawnext))
		w.WriteHeader(http.StatusSeeOther)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var sc, err = r.Cookie("sess")
		if err != nil {
			logger.Error("error obteniendo cookie de sesión: %s", err.Error())
			gotoLogin(w, r.URL.Path)
			return
		}
		sess, err := s.getSession(sc.Value)
		if err != nil {
			logger.Error("error accediendo a página que requiere autenticación: %s", err.Error())
			gotoLogin(w, r.URL.Path)
			return
		}
		newContext := context.WithValue(r.Context(), sessionContextKey("sess"), sess)
		r = r.WithContext(newContext)
		h(w, r)
	}
}
