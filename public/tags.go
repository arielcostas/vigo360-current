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
	"github.com/gorilla/mux"
)

func TagsIdPage(w http.ResponseWriter, r *http.Request) {
	req_tagid := mux.Vars(r)["tagid"]

	tag := Tag{}
	err := db.QueryRowx("SELECT id,nombre FROM tags WHERE id=?;", req_tagid).StructScan(&tag)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Notice("[tagsid]: tried to access unexistent tag %s", req_tagid)
		NotFoundHandler(w, r)
		return
	} else if err != nil {
		logger.Error("[tagsid]: error fetching tag info from database: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	posts := []ResumenPost{}
	err = db.Select(&posts, `SELECT pp.id, DATE_FORMAT(pp.fecha_publicacion, '%d %b. %Y') as fecha_publicacion, pp.alt_portada, pp.titulo, autores.nombre FROM publicaciones_tags LEFT JOIN PublicacionesPublicas pp ON publicaciones_tags.publicacion_id = pp.id LEFT JOIN autores ON pp.autor_id = autores.id WHERE tag_id = ? ORDER BY pp.fecha_publicacion DESC;`, req_tagid)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Notice("[tagsid]: no posts found for tag %s", req_tagid)
		NotFoundHandler(w, r)
		return
	} else if err != nil {
		logger.Error("[tagsid]: error fetching posts from database: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	t.ExecuteTemplate(w, "tags-id.html", struct {
		Tag   Tag
		Posts []ResumenPost
		Meta  common.PageMeta
	}{
		Tag:   tag,
		Posts: posts,
		Meta: common.PageMeta{
			Titulo:      tag.Nombre,
			Descripcion: "Publicaciones en Vigo360 sobre " + tag.Nombre,
			Canonica:    FullCanonica("/tags/" + req_tagid),
		},
	})
}
