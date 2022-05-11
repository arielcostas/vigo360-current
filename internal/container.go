/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/model"
)

// Un container incluye los repositorios para todos los tipos a los que va a acceder el servidor
type Container struct {
	publicacion model.PublicacionStore
	autor       model.AutorStore
	trabajo     model.TrabajoStore
	tag         model.TagStore
	serie       model.SerieStore
}

func NewMysqlContainer(db *sqlx.DB) *Container {
	return &Container{
		publicacion: model.NewPublicacionStore(db),
		autor:       model.NewAutorStore(db),
		trabajo:     model.NewTrabajoStore(db),
		tag:         model.NewTagStore(db),
		serie:       model.NewSerieStore(db),
	}
}
