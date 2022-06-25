/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package repository

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

type MysqlAvisoStore struct {
	db *sqlx.DB
}

func NewMysqlAvisoStore(db *sqlx.DB) *MysqlAvisoStore {
	return &MysqlAvisoStore{
		db: db,
	}
}

func (s *MysqlAvisoStore) Listar() ([]models.Aviso, error) {
	var avisos []models.Aviso
	err := s.db.Select(&avisos, "SELECT fecha_creacion, titulo, contenido FROM avisos ORDER BY fecha_creacion DESC")

	if err != nil {
		return []models.Aviso{}, err
	}

	return avisos, nil
}

func (s *MysqlAvisoStore) ListarRecientes() ([]models.Aviso, error) {
	listado, err := s.Listar()
	return listado[0:4], err
}
