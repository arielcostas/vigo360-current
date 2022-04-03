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

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/thanhpk/randstr"
)

type AdminLoginParams struct {
	LoginError  bool
	PrefillName string
}

type LoginRow struct {
	Id         string
	Nombre     string
	Contraseña string
}

func viewLogin(w http.ResponseWriter, r *http.Request) *appError {
	verifyLogin(w, r)
	err := t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{})
	if err != nil {
		return newTemplateRenderingAppError(err)
	}
	return nil
}

func doLogin(w http.ResponseWriter, r *http.Request) *appError {
	verifyLogin(w, r)

	if err := r.ParseForm(); err != nil {
		return &appError{Error: err, Message: "error parsing form",
			Response: "Hubo un error procesando el inicio de sesión", Status: 500}
	}

	param_userid := r.PostFormValue("userid")
	param_password := r.PostFormValue("password")

	row := LoginRow{}

	if param_userid == "" || param_password == "" {
		err := t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{
			PrefillName: param_userid,
			LoginError:  true,
		})
		if err != nil {
			return &appError{Error: err, Message: "error rendering template",
				Response: "Alguno de los datos introducidos no son correctos.", Status: 400}
		}
		return nil
	}

	// Fetch from database. If error is no user found, show the error document. If a different error is thrown, show 500
	if err := db.QueryRowx("SELECT id, nombre, contraseña FROM autores WHERE id=?;", param_userid).StructScan(&row); errors.Is(err, sql.ErrNoRows) {
		logger.Error("[/login]: no user matches " + param_userid)
		e2 := t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{
			PrefillName: param_userid,
			LoginError:  true,
		})
		if e2 != nil {
			return &appError{Error: e2, Message: "error rendering template",
				Response: "Alguno de los datos introducidos no es correcto.", Status: 400}
		}
	} else if err != nil {
		return &appError{Error: err, Message: "error fetching user",
			Response: "Hubo un error inesperado.", Status: 500}
	}

	pass := ValidatePassword(param_password, row.Contraseña)

	if !pass {
		e2 := t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{
			PrefillName: param_userid,
			LoginError:  true,
		})
		if e2 != nil {
			return &appError{Error: e2, Message: "error rendering template",
				Response: "Alguno de los datos introducidos no es correcto.", Status: 400}
		}
		return nil
	}

	token := randstr.String(20)

	if _, err := db.Exec("INSERT INTO sesiones VALUES (?, NOW(), false, ?)", token, param_userid); err != nil {
		return &appError{Error: err, Message: "error persisting session token to database",
			Response: "Error guardando datos.", Status: 500}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    token,
		Path:     "/admin",
		Domain:   r.URL.Host,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	defer w.WriteHeader(303)
	defer w.Header().Add("Location", "/admin/dashboard")
	return nil
}
