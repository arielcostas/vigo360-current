/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"html/template"
	"net/http"

	"vigo360.es/new/internal/logger"
)

type appError struct {
	// Error whose message is extracted
	Error error
	// Message logged
	Message string
	// Response to give to the user
	Response string
	// HTTP Reply status code
	Status int
}

var statusToType map[int]string = map[int]string{
	400: "Petición inválida",
	403: "Acceso denegado",
	404: "Página no encontrada",
	500: "Error interno del servidor",
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		var rid = r.Context().Value("rid").(string)
		logger.Error("[%s] `%s` %s: %s", rid, r.URL.Path, err.Message, err.Error.Error())

		w.Header().Add("Vigo360-RID", rid)
		w.WriteHeader(err.Status)

		var tipo string
		if t, ok := statusToType[err.Status]; ok {
			tipo = t
		}

		errorDocument.Execute(w, errorDocumentParams{
			Tipo:        tipo,
			Codigo:      err.Status,
			Explicacion: err.Response,
			Rid:         rid,
		})

	}
}

type errorDocumentParams struct {
	Tipo        string
	Codigo      int
	Explicacion string
	Rid         string
}

var errorDocument = template.Must(template.New("errordoc").Parse(
	`<!DOCTYPE html>
	<html>
	<head>
	<title>{{ .Tipo }} - Vigo360</title>
	</head>
	<body>
	<h1>{{ .Codigo }} {{ .Tipo }}</h1>
	<p>Hubo un error intentando mostrar esta página.</p>
	<strong>{{ .Explicacion }}</strong>
	<hr />
	<p>Si crees que se trata de un error, <a href="mailto:contacto@vigo360.es">contacta con el equipo vigo360</a> y provee el código: <b>{{ .Rid }}</b></p>
	<a href="/">Volver al inicio</a>
	</body>
	</html>`,
))
