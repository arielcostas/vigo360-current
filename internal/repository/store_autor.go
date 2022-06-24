/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package repository

import "vigo360.es/new/internal/models"

type AutorStore interface {
	Listar() ([]models.Autor, error)
	Obtener(string) (models.Autor, error)
	Buscar(string) ([]models.Autor, error)
}
