// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package repository

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

type TrabajoStore interface {
	Listar() (models.Trabajos, error)
	ListarPorAutor(string) (models.Trabajos, error)
	ObtenerPorId(string, bool) (models.Trabajo, error)
}

type MysqlTrabajoStore struct {
	db *sqlx.DB
}
