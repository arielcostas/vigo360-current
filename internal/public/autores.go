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
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/model"
	"vigo360.es/new/internal/templates"
)

type AutoresParams struct {
	Autores []model.Autor
	Meta    PageMeta
}

type AutoresIdParams struct {
	Autor    model.Autor
	Posts    model.Publicaciones
	Trabajos model.Trabajos
	Meta     PageMeta
}

func AutoresIdPage(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db         = database.GetDB()
		as         = model.NewAutorStore(db)
		ps         = model.NewPublicacionStore(db)
		ts         = model.NewTrabajoStore(db)
		req_author = mux.Vars(r)["id"]
	)
	var autor model.Autor
	autor, err := as.Obtener(req_author)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "autor not found", Response: "Autor no encontrado", Status: 404}
		}
		return &appError{Error: err, Message: "unexpected error fetching autor", Response: "Error recuperando datos", Status: 500}
	}

	publicaciones, err := ps.ListarPorAutor(autor.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching posts for autor", Response: "Error recuperando datos", Status: 500}
	}

	trabajos, err := ts.ListarPorAutor(autor.Id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching trabajos for autor", Response: "Error recuperando datos", Status: 500}
	}

	err = templates.Render(w, "autores-id.html", AutoresIdParams{
		Autor:    autor,
		Posts:    publicaciones,
		Trabajos: trabajos,
		Meta: PageMeta{
			Titulo:      autor.Nombre,
			Descripcion: autor.Biografia,
			Canonica:    FullCanonica("/autores/" + autor.Id),
		},
	})

	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Error mostrando la página", Status: 500}
	}

	return nil
}

func AutoresPage(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db = database.GetDB()
		as = model.NewAutorStore(db)
	)
	autores, err := as.Listar()
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching autores", Response: "Error recuperando datos", Status: 500}
	}

	err = templates.Render(w, "autores.html", AutoresParams{
		Autores: autores,
		Meta: PageMeta{
			Titulo:      "Autores",
			Descripcion: "Conoce a los autores y colaboradores de Vigo360.",
			Canonica:    FullCanonica("/autores"),
		},
	})
	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "error mostrando la página", Status: 500}
	}

	return nil
}
