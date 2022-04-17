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
	"net/http"

	"github.com/gorilla/mux"
)

type Autor struct {
	Id         string
	Nombre     string
	Email      string
	Rol        string
	Biografia  string
	Web_url    string
	Web_titulo string
}

type AutoresParams struct {
	Autores []Autor
	Meta    PageMeta
}

type AutoresIdParams struct {
	Autor    Autor
	Posts    []ResumenPost
	Trabajos []ResumenPost
	Meta     PageMeta
}

func AutoresIdPage(w http.ResponseWriter, r *http.Request) *appError {
	req_author := mux.Vars(r)["id"]
	var autor Autor
	err := db.QueryRowx("SELECT id, nombre, email, rol, biografia, web_url, web_titulo FROM autores WHERE id=?", req_author).StructScan(&autor)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "autor not found", Response: "Autor no encontrado", Status: 404}
		}
		return &appError{Error: err, Message: "unexpected error fetching autor", Response: "Error recuperando datos", Status: 500}
	}

	var publicaciones = make([]ResumenPost, 0)
	err = db.Select(&publicaciones, `SELECT id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, resumen FROM PublicacionesPublicas pp WHERE autor_id = ? ORDER BY pp.fecha_publicacion DESC;`, req_author)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching posts for autor", Response: "Error recuperando datos", Status: 500}
	}

	var trabajos = make([]ResumenPost, 0)
	err = db.Select(&trabajos, `SELECT id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, resumen FROM TrabajosPublicos WHERE autor_id = ? ORDER BY TrabajosPublicos.fecha_publicacion DESC;`, req_author)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching trabajos for autor", Response: "Error recuperando datos", Status: 500}
	}

	var output bytes.Buffer
	err = t.ExecuteTemplate(&output, "autores-id.html", AutoresIdParams{
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

	w.Write(output.Bytes())
	return nil
}

func AutoresPage(w http.ResponseWriter, r *http.Request) *appError {
	autores := []Autor{}
	err := db.Select(&autores, `SELECT id, nombre, rol, biografia FROM autores;`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching autores", Response: "Error recuperando datos", Status: 500}
	}

	var output bytes.Buffer
	err = t.ExecuteTemplate(&output, "autores.html", AutoresParams{
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

	w.Write(output.Bytes())
	return nil
}
