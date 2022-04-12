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
	"net/url"
	"os"
	"time"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/model"
)

type AtomEntry struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string

	Titulo       string
	Resumen      string
	Autor_id     string
	Autor_nombre string
	Autor_email  string
	Tag_id       sql.NullString
	Raw_tags     sql.NullString
	Tags         []string
}

// DEPRECATED
// TODO: Get rid of this
type FeedParams struct {
	BaseURL      string
	Id           string
	Nombre       string
	LastUpdate   string
	GeneratorURI string
	Entries      []AtomEntry
}

type AtomParams struct {
	Dominio    string
	Path       string
	Titulo     string
	Subtitulo  string
	LastUpdate string
	Entries    model.Publicaciones
}

func PostsAtomFeed(w http.ResponseWriter, r *http.Request) *appError {
	ps := model.NewPublicacionStore(database.GetDB())
	pp, err := ps.Listar()
	if err != nil {
		return &appError{Error: err, Message: "error obtaining posts", Response: "Error obteniendo datos", Status: 500}
	}
	pp = pp.FiltrarPublicas()

	lastUpdate, err := pp.ObtenerUltimaActualizacion()
	if err != nil {
		return &appError{Error: err, Message: "error parsing date", Response: "Error obteniendo datos", Status: 500}
	}

	var result bytes.Buffer
	err = t.ExecuteTemplate(&result, "atom.xml", AtomParams{
		Dominio:    os.Getenv("DOMAIN"),
		Path:       "/atom.xml",
		Titulo:     "Publicaciones",
		Subtitulo:  "Últimas publicaciones en el sitio web de Vigo360",
		LastUpdate: lastUpdate.Format(time.RFC3339),
		Entries:    pp,
	})
	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Error produciendo feed", Status: 500}
	}
	w.Write(result.Bytes())
	return nil
}

func TrabajosAtomFeed(w http.ResponseWriter, r *http.Request) {
	trabajos := []AtomEntry{}
	err := db.Select(&trabajos, `SELECT trabajos.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email FROM TrabajosPublicos trabajos LEFT JOIN autores ON trabajos.autor_id = autores.id`)

	// An unexpected error
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logger.Warning("[atom] unexpected error selecting trabajos: %s", err.Error())
	}

	writeFeed(w, r, "trabajos-atom.xml", trabajos, "Trabajos", "")
}

func TagsAtomFeed(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db = database.GetDB()
		ts = model.NewTagStore(db)
		ps = model.NewPublicacionStore(db)
	)

	tagid := mux.Vars(r)["tagid"]
	tag, err := ts.Obtener(tagid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "tried generating feed for non-existent tag", Response: "Tag no encontrada", Status: 404}
		}
		return &appError{Error: err, Message: "error getting tag", Response: "Error recuperando datos", Status: 500}
	}
	pp, err := ps.ListarPorTag(tag.Id)
	if err != nil {
		return &appError{Error: err, Message: "error obtaining posts", Response: "Error obteniendo datos", Status: 500}
	}
	pp = pp.FiltrarPublicas()

	lastUpdate, err := pp.ObtenerUltimaActualizacion()
	if err != nil {
		return &appError{Error: err, Message: "error parsing date", Response: "Error obteniendo datos", Status: 500}
	}

	var result bytes.Buffer
	err = t.ExecuteTemplate(&result, "atom.xml", AtomParams{
		Dominio:    os.Getenv("DOMAIN"),
		Path:       "/tags/" + tag.Id + "/atom.xml",
		Titulo:     tag.Nombre,
		Subtitulo:  "Últimas publicaciones con la etiqueta " + tag.Nombre,
		LastUpdate: lastUpdate.Format(time.RFC3339),
		Entries:    pp,
	})
	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Error produciendo feed", Status: 500}
	}
	w.Write(result.Bytes())
	return nil

}

func AutorAtomFeed(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db = database.GetDB()
		as = model.NewAutorStore(db)
		ps = model.NewPublicacionStore(db)
	)

	autorid := mux.Vars(r)["autorid"]
	var autor, err = as.ObtenerBasico(autorid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "tried generating feed for nonexistent author", Response: "Autor no encontrado", Status: 404}
		}
		return &appError{Error: err, Message: "error checking author exists", Response: "Error obteniendo datos", Status: 500}
	}

	pp, err := ps.ListarPorAutor(autorid)
	if err != nil {
		return &appError{Error: err, Message: "error obtaining posts", Response: "Error obteniendo datos", Status: 500}
	}
	pp = pp.FiltrarPublicas()

	lastUpdate, err := pp.ObtenerUltimaActualizacion()
	if err != nil {
		return &appError{Error: err, Message: "error parsing date", Response: "Error obteniendo datos", Status: 500}
	}

	var result bytes.Buffer
	err = t.ExecuteTemplate(&result, "atom.xml", AtomParams{
		Dominio:    os.Getenv("DOMAIN"),
		Path:       "/autores/" + autorid + "/atom.xml",
		Titulo:     autor.Nombre,
		Subtitulo:  "Últimas publicaciones escritas por " + autor.Nombre,
		LastUpdate: lastUpdate.Format(time.RFC3339),
		Entries:    pp,
	})
	if err != nil {
		return &appError{Error: err, Message: "error rendering template", Response: "Error produciendo feed", Status: 500}
	}
	w.Write(result.Bytes())
	return nil
}

func writeFeed(w http.ResponseWriter, r *http.Request, feedName string, items []AtomEntry, nombre string, id string) {
	// TODO: Refactor line above
	var lastUpdate time.Time

	for i := 0; i < len(items); i++ {
		p := &items[i]

		t, err := time.Parse("2006-01-02 15:04:05", p.Fecha_actualizacion)
		if err != nil {
			logger.Error("unexpected error parsing fecha_actualizacion: %s", err.Error())
			InternalServerErrorHandler(w, r)
		}

		if lastUpdate.Before(t) {
			lastUpdate = t
		}

		p.Id = url.PathEscape(p.Id)
	}

	w.Header().Add("Content-Type", "application/atom+xml; charset=utf-8")
	err := t.ExecuteTemplate(w, feedName, &FeedParams{
		BaseURL:      os.Getenv("DOMAIN"),
		LastUpdate:   lastUpdate.Format(time.RFC3339),
		Entries:      items,
		Nombre:       nombre,
		Id:           id,
		GeneratorURI: os.Getenv("SOURCE_URL"),
	})

	if err != nil {
		logger.Error("unexpected error rendering feed %s: %s", feedName, err.Error())
		InternalServerErrorHandler(w, r)
	}
}
