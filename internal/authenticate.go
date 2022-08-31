// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"database/sql"
	"errors"
	"time"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/models"
)

// Verifica la validez de una sesión y devuelve si es válida o no, y un error explicando que falló
func (s *Server) getSession(token string) (models.Session, error) {
	var session models.Session
	// TODO: Refactor esto
	var db = database.GetDB()
	err := db.QueryRowx("SELECT sessid as id, iniciada, id as autor_id, nombre as autor_nombre, rol as autor_rol FROM sesiones LEFT JOIN autores ON sesiones.autor_id = autores.id WHERE sessid = ? AND revocada = false;", token).StructScan(&session)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return models.Session{}, models.ErrInvalidSession
		}
		return models.Session{}, err
	}

	hora, err := time.Parse("2006-01-02 15:04:05", session.Iniciada)
	if err != nil {
		return models.Session{}, err
	}
	if time.Since(hora).Hours() > 6 {
		_, err = db.Exec("UPDATE sesiones SET revocada=true WHERE sessid=?", session.Id)
		if err != nil {
			return models.Session{}, err
		}
		return models.Session{}, models.ErrExpiredSession
	}

	session.Permisos = make(map[string]bool)
	perms, err := db.Queryx("SELECT permiso_id FROM permisos_usuarios WHERE autor_id = ?;", session.Autor_id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.Session{}, models.ErrUnablePermissions
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
