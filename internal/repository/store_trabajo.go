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

type TrabajoStore interface {
	Listar() (models.Trabajos, error)
	ListarPorAutor(string) (models.Trabajos, error)
	ObtenerPorId(string, bool) (models.Trabajo, error)
}

type MysqlTrabajoStore struct {
	db *sqlx.DB
}
