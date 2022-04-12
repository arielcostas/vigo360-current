/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

import (
	"github.com/jmoiron/sqlx"
)

type TagStore struct {
	db *sqlx.DB
}

func NewTagStore(db *sqlx.DB) TagStore {
	return TagStore{
		db: db,
	}
}

func (s *TagStore) Listar() (Publicaciones, error) {
	panic("por implementar")
}

func (s *TagStore) Obtener(tag_id string) (Tag, error) {
	var tag Tag
	var row = s.db.QueryRow(`SELECT id, nombre FROM tags WHERE id=?`, tag_id)
	var err = row.Scan(&tag.Id, &tag.Nombre)

	if err != nil {
		return Tag{}, err
	}

	return tag, nil
}
