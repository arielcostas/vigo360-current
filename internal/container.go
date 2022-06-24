/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/repository"
)

// Un container incluye los repositorios para todos los tipos a los que va a acceder el servidor
type Container struct {
	autor       repository.AutorStore
	publicacion repository.PublicacionStore
	serie       repository.SerieStore
	tag         repository.TagStore
	trabajo     repository.TrabajoStore
	comentario  repository.ComentarioStore
}

func NewMysqlContainer(db *sqlx.DB) *Container {
	return &Container{
		autor:       repository.NewMysqlAutorStore(db),
		publicacion: repository.NewMysqlPublicacionStore(db),
		serie:       repository.NewMysqlSerieStore(db),
		tag:         repository.NewMysqlTagStore(db),
		trabajo:     repository.NewMysqlTrabajoStore(db),
		comentario:  repository.NewMysqlComentarioStore(db),
	}
}
