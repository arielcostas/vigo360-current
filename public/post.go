/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"bytes"
	"database/sql"
	"errors"
	"html/template"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
)

func PostPage(w http.ResponseWriter, r *http.Request) {
	req_post_id := mux.Vars(r)["postid"]
	query := `SELECT pp.id, alt_portada, titulo, resumen, contenido, DATE_FORMAT(pp.fecha_publicacion, '%d %b.') as fecha_publicacion, DATE_FORMAT(pp.fecha_actualizacion, '%e %b.') as fecha_actualizacion, autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol, serie_id as serie FROM PublicacionesPublicas pp LEFT JOIN autores on pp.autor_id = autores.id WHERE pp.id = ? ORDER BY pp.fecha_publicacion DESC;`

	post := FullPost{}
	err := db.QueryRowx(query, req_post_id).StructScan(&post)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warning("[post] could not find post with that id")
			NotFoundHandler(w, r)
			return
		}
		logger.Error("[post] unexpected error fetching post from database: %s", err.Error())
	}

	// Fetch series
	var serie Serie
	if post.Serie.Valid {
		serie = Serie{}
		err := db.QueryRowx(`SELECT titulo FROM series WHERE id = ?;`, post.Serie.String).Scan(&serie.Titulo)
		if err != nil {
			logger.Warning("[post] error fetching serie for post %s: %s", post.Id, err.Error())
			InternalServerErrorHandler(w, r)
			return
		}

		err = db.Select(&serie.Articulos, `SELECT id, titulo, serie_posicion FROM PublicacionesPublicas WHERE serie_id=? ORDER BY serie_posicion ASC, titulo ASC`, post.Serie.String)
		if err != nil {
			logger.Warning("[post] error fetching serie for post %s: %s", post.Id, err.Error())
			InternalServerErrorHandler(w, r)
			return
		}
	}

	// TODO: Do this as a template function
	// Result is in markdown, convert to HTML
	var buf bytes.Buffer
	err = common.Parser.Convert([]byte(post.ContenidoRaw), &buf)
	if err != nil {
		logger.Error("[post] error converting markdown to HTML: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	post.Contenido = template.HTML(buf.Bytes())

	err = t.ExecuteTemplate(w, "post.html", struct {
		Post  FullPost
		Meta  common.PageMeta
		Serie Serie
	}{
		Serie: serie,
		Post:  post,
		Meta: common.PageMeta{
			Titulo:      post.Titulo,
			Descripcion: post.Resumen,
			Canonica:    FullCanonica("/post/" + post.Id),
			Miniatura:   FullCanonica("/static/thumb/" + post.Id + ".jpg"),
		},
	})
	if err != nil {
		logger.Error("[autores] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
}
