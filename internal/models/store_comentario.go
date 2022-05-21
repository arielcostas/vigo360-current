/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package models

type ComentarioStore interface {
	// Lista los comentarios públicos para un artículo en forma de lista
	ListarPublicos(publicacion_id string) ([]Comentario, error)
	// Lista los comentarios con un estado específico
	ListarPorEstado(EstadoComentario) ([]Comentario, error)

	GuardarComentario(Comentario) error
	// Cambia el estado de PENDIENTE a APROBADO
	Aprobar(comentario_id string) error
	// Cambia el estado de PENDIENTE a RECHAZADO
	Rechazar(comentario_id string) error
}
