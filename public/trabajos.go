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

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
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

func TrabajoPage(w http.ResponseWriter, r *http.Request) {
	req_paper_id := mux.Vars(r)["paperid"]
	query := `SELECT trabajos.id, alt_portada, titulo, resumen, contenido, 
	DATE_FORMAT(trabajos.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	DATE_FORMAT(trabajos.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol
FROM trabajos 
LEFT JOIN autores on trabajos.autor_id = autores.id 
WHERE trabajos.id = ?;`

	trabajo := Trabajo{}
	err := db.QueryRowx(query, req_paper_id).StructScan(&trabajo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warning("[trabajos] could not find post with that id")
			NotFoundHandler(w, r)
			return
		}

		logger.Error("[trabajos] unexpected error fetching post from database: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	adjuntos := []Adjunto{}
	err = db.Select(&adjuntos, "SELECT nombre_archivo, titulo FROM adjuntos WHERE trabajo_id = ?;", trabajo.Id)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("[trabajos] fetching post attachments: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	t.ExecuteTemplate(w, "trabajos-id.html", struct {
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
}

func TrabajoListPage(w http.ResponseWriter, r *http.Request) {
	trabajos := []ResumenPost{}
	err := db.Select(&trabajos, `SELECT trabajos.id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, autores.nombre FROM trabajos LEFT JOIN autores on trabajos.autor_id = autores.id WHERE trabajos.fecha_publicacion < NOW() ORDER BY trabajos.fecha_publicacion DESC`)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Error("[trabajos] error getting trabajos list: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	t.ExecuteTemplate(w, "trabajos.html", struct {
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
}
