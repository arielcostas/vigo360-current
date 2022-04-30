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
	"strconv"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/model"
	"vigo360.es/new/internal/templates"
)

type indexParams struct {
	CurrentPage int
	PageCount   int
	Posts       model.Publicaciones
	Meta        PageMeta
}

func indexPage(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db = database.GetDB()
		ps = model.NewPublicacionStore(db)
	)

	posts, err := ps.Listar()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching posts", Response: "Error recuperando datos", Status: 500}
	}

	posts = posts.FiltrarPublicas()

	/* Paginación */
	var pagina = 1

	var queryPage = r.URL.Query().Get("page")
	if queryPage != "" {
		o, err := strconv.Atoi(queryPage)
		if err != nil {
			return &appError{err, "queryPage is not an integer", "La página solicitada no es válida", 400}
		}
		pagina = o
	}

	var limite = getMinimo(pagina*9, len(posts)-1)
	var inicio = pagina*9 - 9
	if pagina > 1 {
		limite++
		inicio++
	}

	if inicio >= len(posts) || inicio < 0 {
		return &appError{ErrInvalidInput, "page requested out of bounds", "Página no encontrada", 404}
	}

	if pagina == 1 {
		limite++
	}

	err = templates.Render(w, "index.html", indexParams{
		CurrentPage: pagina,
		PageCount:   (len(posts) / 9) + 1,
		Posts:       posts[inicio:limite],
		Meta: PageMeta{
			Titulo:      "Inicio",
			Descripcion: "Vigo360 es un proyecto dedicado a estudiar varios aspectos de la ciudad de Vigo (España) y su área de influencia, centrándose en la toponimia y el transporte.",
			Canonica:    FullCanonica("/"),
		},
	})
	if err != nil {
		// TODO estandarizar esto
		return &appError{err, "error rendering template", "Error mostrando la página", 500}
	}

	return nil
}
