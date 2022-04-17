/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"bytes"
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("page not found")

func nodbPageError(err error) *appError {
	return &appError{Error: err, Message: "error rendering template", Response: "Error mostrando la página", Status: 500}
}

var nodbPageMeta = map[string]PageMeta{
	"siguenos": {
		Titulo:      "Síguenos",
		Descripcion: "Información sobre cómo seguir a Vigo360, y enterarse de sus últimas publicaciones y novedades.",
		Canonica:    FullCanonica("/siguenos"),
	},
	"legal": {
		Titulo:      "Licencias",
		Descripcion: "Información legal relativa a Vigo360, desde licencias de uso libre hasta la política de privacidad.",
		Canonica:    FullCanonica("/licencia"),
	},
	"contacto": {
		Titulo:      "Contacto",
		Descripcion: "Si necesitases contactar con Vigo360, aquí encontrarás cómo hacerlo.",
		Canonica:    FullCanonica("/contacto"),
	},
}

func NoDbPage(w http.ResponseWriter, r *http.Request) *appError {
	var (
		page   = r.URL.Path[1:]
		meta   PageMeta
		output bytes.Buffer
	)

	if pm, ok := nodbPageMeta[page]; ok {
		meta = pm
	} else {
		return &appError{Error: ErrNotFound, Message: "page not found", Response: "La página solicitada no se ha encontrado.", Status: 404}
	}

	err := t.ExecuteTemplate(&output, page+".html", NoPageData{
		Meta: meta,
	})
	if err != nil {
		return nodbPageError(err)
	}

	w.Write(output.Bytes())
	return nil
}
