/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package models

type EstadoComentario int

const (
	ESTADO_PENDIENTE EstadoComentario = 1
	ESTADO_APROBADO  EstadoComentario = 2
	ESTADO_RECHAZADO EstadoComentario = 3
)

type Comentario struct {
	Id             string
	Publicacion_id string
	Padre_id       string

	Nombre         string
	Es_autor       bool
	Autor_original bool
	Contenido      string

	Fecha_creacion   string
	Fecha_moderacion string
	Estado           EstadoComentario
	Moderador        string
}
