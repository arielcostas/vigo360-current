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

type Trabajo struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string
	Alt_portada         string
	Titulo              string
	Resumen             string
	Contenido           string
	Autor_id            string
	Autor_nombre        string
	Autor_rol           string
	Autor_biografia     string
}

type Adjunto struct {
	Nombre_archivo string
	Titulo         string
}

func listTrabajos(w http.ResponseWriter, r *http.Request) *appError {
	trabajos := []ResumenPost{}
	err := db.Select(&trabajos, `SELECT trabajos.id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, autores.nombre FROM TrabajosPublicos trabajos LEFT JOIN autores on trabajos.autor_id = autores.id ORDER BY trabajos.fecha_publicacion DESC`)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error listing trabajos", Response: "Error recuperando datos", Status: 500}
	}

	var output bytes.Buffer
	err = t.ExecuteTemplate(&output, "trabajos.html", struct {
		Trabajos []ResumenPost
		Meta     PageMeta
	}{
		Trabajos: trabajos,
		Meta: PageMeta{
			Titulo:      "Trabajos",
			Descripcion: "Trabajos originales e interesantes publicados por los autores de Vigo360.",
			Canonica:    FullCanonica("/trabajos"),
		},
	})

	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Hubo un error mostrando la página", Status: 500}
	}

	w.Write(output.Bytes())
	return nil
}

func viewTrabajo(w http.ResponseWriter, r *http.Request) *appError {
	trabajoid := mux.Vars(r)["trabajoid"]

	var trabajo Trabajo
	if nt, err := GetFullTrabajo(trabajoid); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "trabajo not found", Response: "El trabajo buscado no se ha encontrado", Status: 404}
		}
		return &appError{Error: err, Message: "error fetching trabajo", Response: "Error recuperando datos", Status: 500}
	} else {
		trabajo = nt
	}

	var adjuntos = make([]Adjunto, 0)
	err := db.Select(&adjuntos, "SELECT nombre_archivo, titulo FROM adjuntos WHERE trabajo_id = ?;", trabajo.Id)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return &appError{Error: err, Message: "error fetching attachments", Response: "Error recuperando datos", Status: 500}
	}

	var output bytes.Buffer
	err = t.ExecuteTemplate(&output, "trabajos-id.html", struct {
		Trabajo  Trabajo
		Adjuntos []Adjunto
		Meta     PageMeta
	}{
		Trabajo:  trabajo,
		Adjuntos: adjuntos,
		Meta: PageMeta{
			Titulo:      trabajo.Titulo,
			Descripcion: trabajo.Resumen,
			Canonica:    FullCanonica("/trabajos/" + trabajo.Id),
			Miniatura:   FullCanonica("/static/thumb/" + trabajo.Id + ".jpg"),
		},
	})

	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Error mostrando la página", Status: 500}
	}

	w.Write(output.Bytes())
	return nil
}
