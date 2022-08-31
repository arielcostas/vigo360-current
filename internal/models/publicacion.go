// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package models

type Publicacion struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string
	Alt_portada         string
	Titulo              string
	Resumen             string
	Contenido           string
	Comentarios         []Comentario

	Serie          Serie
	Serie_posicion int
	Autor          Autor
	Tags           []Tag
}
