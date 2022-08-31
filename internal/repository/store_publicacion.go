// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

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
