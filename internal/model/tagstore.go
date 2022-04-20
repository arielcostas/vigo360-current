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

func (s *TagStore) Listar() ([]Tag, error) {
	var tags = make([]Tag, 0)
	var rows, err = s.db.Query(`SELECT id, nombre FROM tags`)
	if err != nil {
		return []Tag{}, err
	}

	for rows.Next() {
		var nt Tag
		err = rows.Scan(&nt.Id, &nt.Nombre)
		if err != nil {
			return []Tag{}, err
		}
		tags = append(tags, nt)
	}

	return tags, nil
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
