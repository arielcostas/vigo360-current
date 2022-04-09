/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/go-playground/validator/v10"
)

func listSeries(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}
	_, err = getSession(sc.Value)
	if err != nil {
		logger.Notice("unauthenticated user tried to access this page")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}

	series := []Serie{}
	err = db.Select(&series, `SELECT series.*, COUNT(publicaciones.id) as articulos FROM series LEFT JOIN publicaciones ON series.id = publicaciones.serie_id GROUP BY series.id;`)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return newDatabaseReadAppError(err, "series")
	}

	err = t.ExecuteTemplate(w, "series.html", struct {
		Series []Serie
	}{
		Series: series,
	})
	if err != nil {
		return &appError{Error: err, Message: "error rendering page",
			Response: "Hubo un error mostrando esta página.", Status: 500}
	}
	return nil
}

type CreateSeriesFormInput struct {
	Titulo string `validate:"required,min=1,max=40"`
}

func createSeries(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}
	_, err = getSession(sc.Value)
	if err != nil {
		logger.Notice("unauthenticated user tried to access this page")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		return &appError{Error: err, Message: "error parsing form",
			Response: "Error recibiendo los datos. Inténtalo de nuevo más tarde.", Status: 500}
	}

	fi := CreateSeriesFormInput{}
	fi.Titulo = r.FormValue("titulo")

	if err := validator.New().Struct(fi); err != nil {
		return &appError{Error: err, Message: "error validating form",
			Response: "Alguno de los datos introducidos no es válido.", Status: 400}
	}

	id := strings.ToLower(strings.TrimSpace(fi.Titulo))
	id = strings.ReplaceAll(id, " ", "-")

	if _, err := db.Exec(`INSERT INTO series VALUES (?, ?)`, id, fi.Titulo); err != nil {
		return &appError{Error: err, Message: "error persisting to database",
			Response: "Hubo un error guardando datos.", Status: 500}
	}

	w.Header().Add("Location", "/admin/series")
	w.WriteHeader(303)
	return nil
}
