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
	publicacion models.PublicacionStore
	autor       models.AutorStore
	trabajo     models.TrabajoStore
	tag         models.TagStore
	serie       models.SerieStore
}

func NewMysqlContainer(db *sqlx.DB) *Container {
	return &Container{
		publicacion: models.NewPublicacionStore(db),
		autor:       models.NewAutorStore(db),
		trabajo:     models.NewTrabajoStore(db),
		tag:         models.NewTagStore(db),
		serie:       models.NewSerieStore(db),
	}
}
