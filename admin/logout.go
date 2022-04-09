/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"net/http"
	"time"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

func logoutPage(w http.ResponseWriter, r *http.Request) *appError {
	sc, err := r.Cookie("sess")
	if err != nil { // User isn't logged in
		http.Redirect(w, r, "/admin/login", 302)
		return nil
	}

	sess, _ := getSession(sc.Value)
	if err := revokeSession(sess.Id); err != nil {
		return &appError{Error: err, Message: "error revoking session with token " + sess.Id,
			Response: "Hubo un error cerrando la sesión", Status: 500}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    "",
		Path:     "/admin",
		Domain:   r.URL.Host,
		Expires:  time.Date(1970, 01, 01, 00, 00, 00, 00, time.UTC),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	logger.Information("revoked session with id %s", sess.Id)
	w.Header().Add("Location", "/admin/login")
	w.WriteHeader(302)
	return nil
}
