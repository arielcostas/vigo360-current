/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package models

import "github.com/jmoiron/sqlx"

type MysqlSerieStore struct {
	db *sqlx.DB
}

func NewMysqlSerieStore(db *sqlx.DB) *MysqlSerieStore {
	return &MysqlSerieStore{
		db: db,
	}
}

func (s *MysqlSerieStore) Listar() (Publicaciones, error) {
	panic("por implementar")
}

func (s *MysqlSerieStore) Obtener(serie_id string) (Serie, error) {
	var serie Serie
	var row = s.db.QueryRow(`SELECT id, titulo FROM series WHERE id=?`, serie_id)
	var err = row.Scan(&serie.Id, &serie.Titulo)

	if err != nil {
		return Serie{}, err
	}

	filas, err := s.db.Query(`SELECT id, titulo, fecha_publicacion FROM publicaciones WHERE serie_id=?`, serie_id)
	if err != nil {
		return Serie{}, err
	}

	for filas.Next() {
		na := Publicacion{}
		filas.Scan(&na.Id, &na.Titulo, &na.Fecha_publicacion)
		serie.Publicaciones = append(serie.Publicaciones, na)
	}

	return serie, nil
}
