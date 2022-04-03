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
)

func viewDashboard(w http.ResponseWriter, r *http.Request) *appError {
	sesion := verifyLogin(w, r)

	avisos := []Aviso{}
	err := db.Select(&avisos, "SELECT DATE_FORMAT(fecha_creacion, '%d %b.') as fecha_creacion, titulo, contenido FROM avisos ORDER BY avisos.fecha_creacion DESC LIMIT 5")

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return newDatabaseReadAppError(err, "avisos")
	}

	posts := []DashboardPost{}
	err = db.Select(&posts, "SELECT publicaciones.id, titulo, DATE_FORMAT(fecha_publicacion, '%d %b.') as fecha_publicacion, resumen, autores.nombre as autor_nombre FROM publicaciones LEFT JOIN autores ON publicaciones.autor_id = autores.id ORDER BY publicaciones.fecha_publicacion DESC LIMIT 5;")

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return newDatabaseReadAppError(err, "posts")
	}

	err = t.ExecuteTemplate(w, "admin-dashboard.html", struct {
		Sesion Sesion
		Avisos []Aviso
		Posts  []DashboardPost
	}{
		Sesion: sesion,
		Avisos: avisos,
		Posts:  posts,
	})
	if err != nil {
		return newTemplateRenderingAppError(err)
	}
	return nil
}
