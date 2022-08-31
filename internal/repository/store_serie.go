// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package repository

import "vigo360.es/new/internal/models"

type SerieStore interface {
	Listar() ([]models.Serie, error)
	Obtener(string) (models.Serie, error)
}
