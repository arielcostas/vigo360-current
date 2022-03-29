/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"net/http"
	"strings"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/go-playground/validator/v10"
)

func ListSeries(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	series := []Serie{}
	err := db.Select(&series, `SELECT series.*, COUNT(publicaciones.id) as articulos FROM series LEFT JOIN publicaciones ON series.id = publicaciones.serie_id GROUP BY series.id;`)
	if err != nil {
		logger.Error("[series]: error fetching series from database: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	err = t.ExecuteTemplate(w, "series.html", struct {
		Series []Serie
	}{
		Series: series,
	})
}

type CreateSeriesFormInput struct {
	Titulo string `validate:"required,min=1.max=40"`
}

func CreateSeries(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	err := r.ParseForm()

	if err != nil {
		logger.Error("[create-series] error parsing form: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}

	fi := CreateSeriesFormInput{}
	fi.Titulo = r.FormValue("titulo")

	err = validator.New().Struct(fi)
	if err != nil {
		logger.Error("[serie] error validating form: %s", err.Error())
		w.WriteHeader(400)
		w.Write([]byte("Error de validación"))
		return
	}

	id := strings.ToLower(strings.TrimSpace(fi.Titulo))
	id = strings.ReplaceAll(id, " ", "-")

	_, err = db.Exec(`INSERT INTO series VALUES (?, ?)`, id, fi.Titulo)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Error guardando nueva serie en la base de datos"))
		logger.Error("[create-series] error saving new series to database: %s", err.Error())
		return
	}

	w.Header().Add("Location", "/admin/series")
	w.WriteHeader(303)
}
