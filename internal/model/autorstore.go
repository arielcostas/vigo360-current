/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

import (
	"github.com/jmoiron/sqlx"
)

type AutorStore struct {
	db *sqlx.DB
}

func NewAutorStore(db *sqlx.DB) AutorStore {
	return AutorStore{
		db: db,
	}
}

func (s *AutorStore) Listar() (Publicaciones, error) {
	panic("por implementar")
}

func (s *AutorStore) ObtenerBasico(autor_id string) (Autor, error) {
	var autor Autor
	var row = s.db.QueryRow(`SELECT id, nombre, email, rol, biografia, web_url, web_titulo FROM autores WHERE id=?`, autor_id)
	var err = row.Scan(&autor.Id, &autor.Nombre, &autor.Email, &autor.Rol, &autor.Biografia, &autor.Web.Url, &autor.Web.Titulo)

	if err != nil {
		return Autor{}, err
	}

	return autor, nil
}
