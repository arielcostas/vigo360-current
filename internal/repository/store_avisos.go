// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package repository

import "vigo360.es/new/internal/models"

type AvisoStore interface {
	// Obtiene todos los avisos
	Listar() ([]models.Aviso, error)
	// Obtiene los 5 avisos más recientes
	ListarRecientes() ([]models.Aviso, error)
}
