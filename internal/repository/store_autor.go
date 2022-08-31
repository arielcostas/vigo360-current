// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package repository

import "vigo360.es/new/internal/models"

type AutorStore interface {
	Listar() ([]models.Autor, error)
	Obtener(string) (models.Autor, error)
	Buscar(string) ([]models.Autor, error)
}
