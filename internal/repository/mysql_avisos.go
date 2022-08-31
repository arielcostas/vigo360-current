// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

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
