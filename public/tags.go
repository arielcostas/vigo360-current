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

	"github.com/gorilla/mux"
)

func listTags(w http.ResponseWriter, r *http.Request) *appError {
	var tags = make([]Tag, 0)
	if err := db.Select(&tags, `SELECT *, (SELECT COUNT(publicacion_id) FROM publicaciones_tags LEFT JOIN publicaciones p2 ON publicaciones_tags.publicacion_id = p2.id WHERE tag_id = tags.id AND p2.fecha_publicacion < NOW()) as publicaciones, (SELECT publicacion_id FROM publicaciones_tags LEFT JOIN publicaciones p2 ON publicaciones_tags.publicacion_id = p2.id WHERE tag_id = tags.id AND p2.fecha_publicacion < NOW() LIMIT 1) as ultima FROM tags HAVING publicaciones > 0;`); err != nil {
		return &appError{Error: err, Message: "error fetching tags", Response: "Hubo un error mostrando esta página.", Status: 500}
	}

	err := t.ExecuteTemplate(w, "tags.html", struct {
		Tags []Tag
		Meta PageMeta
	}{
		Tags: tags,
		Meta: PageMeta{
			Titulo:      "Tags",
			Descripcion: "Las diversas tags en las que se categorizan los artículos de Vigo360",
			Canonica:    FullCanonica("/tags"),
		},
	})

	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Hubo un error mostrando esta página.", Status: 500}
	}

	return nil
}

func viewTag(w http.ResponseWriter, r *http.Request) *appError {
	req_tagid := mux.Vars(r)["tagid"]

	tag := Tag{}
	err := db.QueryRowx("SELECT id,nombre FROM tags WHERE id=?;", req_tagid).StructScan(&tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "tag not found", Response: "La página buscada no se pudo encontrar", Status: 404}
		}
		return &appError{Error: err, Message: "unexpected error fetching tag", Response: "Error recuperando datos", Status: 500}
	}

	posts := []ResumenPost{}
	err = db.Select(&posts, `SELECT pp.id, DATE_FORMAT(pp.fecha_publicacion, '%d %b. %Y') as fecha_publicacion, pp.alt_portada, pp.titulo, autores.nombre FROM publicaciones_tags RIGHT JOIN PublicacionesPublicas pp ON publicaciones_tags.publicacion_id = pp.id LEFT JOIN autores ON pp.autor_id = autores.id WHERE tag_id = ? ORDER BY pp.fecha_publicacion DESC;`, req_tagid)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching posts for tag", Response: "Error recuperando datos", Status: 500}
	}

	err = t.ExecuteTemplate(w, "tags-id.html", struct {
		Tag   Tag
		Posts []ResumenPost
		Meta  PageMeta
	}{
		Tag:   tag,
		Posts: posts,
		Meta: PageMeta{
			Titulo:      tag.Nombre,
			Keywords:    tag.Nombre,
			Descripcion: "Publicaciones en Vigo360 sobre " + tag.Nombre,
			Canonica:    FullCanonica("/tags/" + req_tagid),
		},
	})

	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Error mostrando la página.", Status: 500}
	}

	return nil
}
