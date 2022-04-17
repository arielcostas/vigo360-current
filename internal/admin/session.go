package admin

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
)

func listSessions(w http.ResponseWriter, r *http.Request) *appError {
	return nil
}

func revokeSession(sessid string) error {
	_, err := db.Exec("UPDATE sesiones SET revocada = 1 WHERE sessid = ?;", sessid)
	if err != nil {
		return err
	}

	return nil
}

/*
	Verifies a login token's validity and returns whether is valid or not
	and an error explaining what went wrong
*/
func getSession(token string) (Session, error) {
	var session Session
	err := db.QueryRowx("SELECT sessid as id, iniciada, id as autor_id, nombre as autor_nombre, rol as autor_rol FROM sesiones LEFT JOIN autores ON sesiones.autor_id = autores.id WHERE sessid = ? AND revocada = false;", token).StructScan(&session)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return Session{}, ErrInvalidSession
		}
		return Session{}, err
	}

	hora, err := time.Parse("2006-01-02 15:04:05", session.Iniciada)
	if err != nil {
		return Session{}, err
	}
	if time.Since(hora).Hours() > 6 {
		_, err = db.Exec("UPDATE sesiones SET revocada=true WHERE sessid=?", session.Id)
		if err != nil {
			return Session{}, err
		}
		return Session{}, ErrExpiredSession
	}

	session.Permisos = make(map[string]bool)
	perms, err := db.Queryx("SELECT permiso_id FROM permisos_usuarios WHERE autor_id = ?;", session.Autor_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return Session{}, ErrUnablePermissions
	}

	for perms.Next() {
		var p string
		err = perms.Scan(&p)
		if err != nil {
			continue
		}
		session.Permisos[p] = true
	}

	return session, nil
}
