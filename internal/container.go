/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/models"
)

// Un container incluye los repositorios para todos los tipos a los que va a acceder el servidor
type Container struct {
	autor       models.AutorStore
	publicacion models.PublicacionStore
	serie       models.SerieStore
	tag         models.TagStore
	trabajo     models.TrabajoStore
}

func NewMysqlContainer(db *sqlx.DB) *Container {
	return &Container{
		autor:       models.NewMysqlAutorStore(db),
		publicacion: models.NewMysqlPublicacionStore(db),
		serie:       models.NewMysqlSerieStore(db),
		tag:         models.NewMysqlTagStore(db),
		trabajo:     models.NewMysqlTrabajoStore(db),
	}
}
