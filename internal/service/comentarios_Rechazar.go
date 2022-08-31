// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package service

func (se *Comentario) Rechazar(comentario_id string, moderador_id string) error {
	return se.cstore.Rechazar(comentario_id, moderador_id)
}
