/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package repository

import "vigo360.es/new/internal/models"

type AvisoStore interface {
	// Obtiene todos los avisos
	Listar() ([]models.Aviso, error)
	// Obtiene los 5 avisos más recientes
	ListarRecientes() ([]models.Aviso, error)
}
