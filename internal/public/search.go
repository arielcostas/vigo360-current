/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"net/http"
	"strings"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/model"
	"vigo360.es/new/internal/templates"
)

type searchPageParams struct {
	Resultados []ResultadoBusqueda
	Termino    string
	Meta       PageMeta
}

type ResultadoBusqueda struct {
	Id      string
	Titulo  string
	Resumen string
	Uri     string
}

func realizarBusqueda(w http.ResponseWriter, r *http.Request) *appError {
	var (
		db         = database.GetDB()
		ps         = model.NewPublicacionStore(db)
		as         = model.NewAutorStore(db)
		resultados = make([]ResultadoBusqueda, 0)
	)

	var termino = r.URL.Query().Get("termino")
	termino = strings.TrimSpace(termino)
	// TODO Gestionar-impedir términos vacíos
	if termino == "" {
		w.Header().Add("Location", "/")
		w.WriteHeader(302)
		return nil
	}

	autores, err := as.Buscar(termino)
	if err != nil {
		return &appError{err, "error searching authors", "Hubo un error realizando la búsqueda", 500}
	}

	for _, autor := range autores {
		resultados = append(resultados, ResultadoBusqueda{
			Id:      autor.Id,
			Titulo:  autor.Nombre,
			Resumen: autor.Biografia,
			Uri:     "/autores/" + autor.Id,
		})
	}

	publicaciones, err := ps.Buscar(termino)
	publicaciones = publicaciones.FiltrarPublicas()
	if err != nil {
		return &appError{err, "error searching posts", "Hubo un error realizando la búsqueda", 500}
	}

	for _, pub := range publicaciones {
		resultados = append(resultados, ResultadoBusqueda{
			Id:      pub.Id,
			Titulo:  pub.Titulo,
			Resumen: pub.Resumen,
			Uri:     "/post/" + pub.Id,
		})
	}

	err = templates.Render(w, "search.html", searchPageParams{
		Resultados: resultados,
		Termino:    termino,
		Meta: PageMeta{
			Titulo:   "Resultados para " + termino,
			Canonica: FullCanonica("/buscar?termino=" + termino),
		},
	})
	if err != nil {
		return &appError{err, "error rendering template", "Hubo un error mostrando la página solicitada", 500}
	}
	return nil
}
