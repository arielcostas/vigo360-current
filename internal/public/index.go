/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"database/sql"
	"errors"
	"net/http"
)

func indexPage(w http.ResponseWriter, r *http.Request) *appError {
	posts := []ResumenPost{}
	err := db.Select(&posts, "SELECT pp.id, DATE_FORMAT(pp.fecha_publicacion, '%d %b. %Y') as fecha_publicacion, pp.alt_portada, pp.titulo, pp.resumen, autores.nombre FROM PublicacionesPublicas pp LEFT JOIN autores on pp.autor_id = autores.id ORDER BY pp.fecha_publicacion DESC;")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching posts", Response: "Error recuperando datos", Status: 500}
	}

	err = t.ExecuteTemplate(w, "index.html", struct {
		Posts []ResumenPost
		Meta  PageMeta
	}{
		Posts: posts,
		Meta: PageMeta{
			Titulo:      "Inicio",
			Descripcion: "Vigo360 es un proyecto dedicado a estudiar varios aspectos de la ciudad de Vigo (España) y su área de influencia, centrándose en la toponimia y el transporte.",
			Canonica:    FullCanonica("/"),
		},
	})

	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Error mostrando la página", Status: 500}
	}

	return nil
}
