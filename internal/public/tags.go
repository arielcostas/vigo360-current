/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"database/sql"
	"errors"
	"math/rand"
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/model"
	"vigo360.es/new/internal/templates"
)

type listTagsParams struct {
	Tags []model.Tag
	Meta PageMeta
}

type viewTagParams struct {
	Tag   model.Tag
	Posts model.Publicaciones
	Meta  PageMeta
}

func listTags(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db = database.GetDB()
		ts = model.NewTagStore(db)
		ps = model.NewPublicacionStore(db)
	)

	var tags, err = ts.Listar()
	for i, t := range tags {
		if t.Publicaciones < 1 {
			tags = append(tags[:i], tags[i+1:]...)
		}
	}

	if err != nil {
		return &appError{err, "error fetching tags", "Hubo un error obteniendo datos.", 500}
	}

	sort.Slice(tags, func(p, q int) bool {
		return tags[p].Nombre < tags[q].Nombre
	})

	/* Obtener última publicación para cada tag evitando duplicaciones */
	var publicacionesUsadas = make(map[string]bool, 0)
	for i, t := range tags {
		nt := t
		var publicacionesConTag []string
		pt, err := ps.ListarPorTag(t.Id)
		pt.FiltrarPublicas()
		for _, p := range pt {
			publicacionesConTag = append(publicacionesConTag, p.Id)
		}
		if err != nil {
			return &appError{err, "error fetching publicaciones with tag " + t.Id, "Hubo un error obteniendo datos.", 500}
		}

		for _, pub := range publicacionesConTag {
			if _, ok := publicacionesUsadas[pub]; !ok {
				publicacionesUsadas[pub] = true
				nt.Ultima = pub
				break
			}
		}

		// Si se diera el caso de que todas están escogidas, poner una al azar
		if nt.Ultima == "" {
			cant := len(publicacionesConTag) - 1
			if cant <= 1 {
				cant = 1
			}
			aleatoria := rand.Intn(cant)
			nt.Ultima = publicacionesConTag[aleatoria]
		}

		tags[i] = nt
	}

	err = templates.Render(w, "tags.html", listTagsParams{
		Tags: tags,
		Meta: PageMeta{
			Titulo:      "Tags",
			Descripcion: "Las diversas tags en las que se categorizan los artículos de Vigo360",
			Canonica:    FullCanonica("/tags"),
		},
	})

	if err != nil {
		return &appError{err, "error rendering template", "Hubo un error mostrando esta página.", 500}
	}

	return nil
}

func viewTag(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db = database.GetDB()
		ps = model.NewPublicacionStore(db)
		ts = model.NewTagStore(db)
	)
	req_tagid := mux.Vars(r)["tagid"]

	tag, err := ts.Obtener(req_tagid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "tag not found", Response: "La página buscada no se pudo encontrar", Status: 404}
		}
		return &appError{Error: err, Message: "unexpected error fetching tag", Response: "Error recuperando datos", Status: 500}
	}

	publicaciones, err := ps.ListarPorTag(req_tagid)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching posts for tag", Response: "Error recuperando datos", Status: 500}
	}

	err = templates.Render(w, "tags-id.html", viewTagParams{
		Tag:   tag,
		Posts: publicaciones.FiltrarPublicas(),
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
