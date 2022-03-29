/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"database/sql"
	"errors"
)

type Sesion struct {
	Sessid   string
	Id       string
	Iniciada string
	Nombre   string
	Rol      string
}

type Aviso struct {
	Fecha_creacion string
	Titulo         string
	Contenido      string
}

type DashboardPost struct {
	Id                string
	Titulo            string
	Resumen           string
	Fecha_publicacion string
	Autor_nombre      string
}

type PostEditar struct {
	Id             string
	Titulo         string
	Resumen        string
	Contenido      string
	Alt_portada    string
	Publicado      bool
	Serie_id       sql.NullString
	Serie_posicion sql.NullInt16
}

type Serie struct {
	Id        string
	Titulo    string
	Articulos int
}

type Tag struct {
	Id           string
	Nombre       string
	Seleccionada bool
}

var InvalidImageFormatError error = errors.New("invalid image MIME type")
