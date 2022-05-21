/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package models

type EstadoComentario int

const (
	ESTADO_PENDIENTE EstadoComentario = iota
	ESTADO_APROBADO  EstadoComentario = iota
	ESTADO_RECHAZADO EstadoComentario = iota
)

type Comentario struct {
	Id             string
	Publicacion_id string
	Padre_id       string

	Nombre         string
	Email_hash     string
	Es_autor       bool
	Autor_original bool
	Contenido      string

	Fecha_creacion   string
	Fecha_moderacion string
	Estado           EstadoComentario
	Moderador        string
}
