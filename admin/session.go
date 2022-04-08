package admin

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

func revokeSession(sessid string) error {
	_, err := db.Exec("UPDATE sesiones SET revocada = 1 WHERE sessid = ?;", sessid)
	return err
}

func gotoLogin(w http.ResponseWriter, r *http.Request) Sesion {
	http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
	return Sesion{}
}

func verifyLogin(w http.ResponseWriter, r *http.Request) Sesion {
	// TODO: Refactor this to return an error or the session AND DO NOT SEND TO LOG IN
	cookie, err := r.Cookie("sess")

	if errors.Is(err, http.ErrNoCookie) && r.URL.Path != "/admin/login" {
		logger.Notice("unauthenticated user tried accessing auth-requiring page %s", r.URL.Path)
		return gotoLogin(w, r)
	}

	if err != nil && r.URL.Path != "/admin/login" {
		logger.Error("error getting session cookie: %s", err.Error())
		return gotoLogin(w, r)
	} else if err != nil {
		return Sesion{}
	}

	user := Sesion{}

	err = db.QueryRowx("SELECT sessid, iniciada, id, nombre, rol FROM sesiones LEFT JOIN autores ON sesiones.autor_id = autores.id WHERE sessid = ? AND revocada = false;", cookie.Value).StructScan(&user)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Warning("error in login verification: %s", err.Error())
	} else if err != nil {
		logger.Error("unexpected error fetching session from database: %s", err.Error())
	}

	if err != nil && r.URL.Path != "/admin/login" {
		return gotoLogin(w, r)
	}

	hora, _ := time.Parse("2006-01-02 15:04:05", user.Iniciada)

	if time.Since(hora).Hours() > 6 {
		db.Exec("UPDATE sesiones SET revocada=true WHERE sessid=?", user.Sessid)
		logger.Warning("session older than 6 hours, revoking automatically")
		if r.URL.Path != "/admin/login" {
			return gotoLogin(w, r)
		} else {
			return Sesion{}
		}
	}

	// Logged in successfully, no sense to log in again
	if r.URL.Path == "/admin/login" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
	}

	// It's not the login page and the user is logged in
	return user
}

func listSessions(w http.ResponseWriter, r *http.Request) *appError {
	return nil
}
