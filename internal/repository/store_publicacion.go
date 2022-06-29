/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package repository

import "vigo360.es/new/internal/models"

type PublicacionStore interface {
	Listar() (models.Publicaciones, error)
	ListarPorAutor(autor_id string) (models.Publicaciones, error)
	ListarPorTag(tag_id string) (models.Publicaciones, error)
	ListarPorSerie(serie_id string) (models.Publicaciones, error)
	Existe(id string) (bool, error)
	ObtenerPorId(id string, requirePublic bool) (models.Publicacion, error)
	Buscar(query string) (models.Publicaciones, error)
}
