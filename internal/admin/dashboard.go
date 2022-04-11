/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"database/sql"
	"errors"
	"net/http"

	"vigo360.es/new/internal/logger"
)

func viewDashboard(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}
	sess, err := getSession(sc.Value)
	if err != nil {
		logger.Notice("unauthenticated user tried to access this page")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}

	avisos := []Aviso{}
	err = db.Select(&avisos, "SELECT DATE_FORMAT(fecha_creacion, '%d %b.') as fecha_creacion, titulo, contenido FROM avisos ORDER BY avisos.fecha_creacion DESC LIMIT 5")

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return newDatabaseReadAppError(err, "avisos")
	}

	posts := []DashboardPost{}
	err = db.Select(&posts, "SELECT publicaciones.id, titulo, DATE_FORMAT(fecha_publicacion, '%d %b.') as fecha_publicacion, resumen, autores.nombre as autor_nombre FROM publicaciones LEFT JOIN autores ON publicaciones.autor_id = autores.id ORDER BY publicaciones.fecha_publicacion DESC LIMIT 5;")

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return newDatabaseReadAppError(err, "posts")
	}

	err = t.ExecuteTemplate(w, "admin-dashboard.html", struct {
		Avisos  []Aviso
		Posts   []DashboardPost
		Session Session
	}{
		Avisos:  avisos,
		Posts:   posts,
		Session: sess,
	})
	if err != nil {
		return newTemplateRenderingAppError(err)
	}
	return nil
}
