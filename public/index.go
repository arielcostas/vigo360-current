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

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	posts := []ResumenPost{}
	err := db.Select(&posts, "SELECT pp.id, DATE_FORMAT(pp.fecha_publicacion, '%d %b. %Y') as fecha_publicacion, pp.alt_portada, pp.titulo, pp.resumen, autores.nombre FROM PublicacionesPublicas pp LEFT JOIN autores on pp.autor_id = autores.id ORDER BY pp.fecha_publicacion DESC;")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("[index]: error fetching posts: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	err = t.ExecuteTemplate(w, "index.html", struct {
		Posts []ResumenPost
		Meta  common.PageMeta
	}{
		Posts: posts,
		Meta: common.PageMeta{
			Titulo:      "Inicio",
			Descripcion: "Vigo360 es un proyecto dedicado a estudiar varios aspectos de la ciudad de Vigo (España) y su área de influencia, centrándose en la toponimia y el transporte.",
			Canonica:    FullCanonica("/"),
		},
	})

	if err != nil {
		logger.Error("[index] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
}
