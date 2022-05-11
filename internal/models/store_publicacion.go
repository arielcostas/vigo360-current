/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package models

type PublicacionStore interface {
	Listar() (Publicaciones, error)
	ListarPorAutor(string) (Publicaciones, error)
	ListarPorTag(string) (Publicaciones, error)
	ListarPorSerie(string) (Publicaciones, error)
	ObtenerPorId(string, bool) (Publicacion, error)
	Buscar(string) (Publicaciones, error)
}
