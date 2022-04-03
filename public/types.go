/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"database/sql"
)

type Autor struct {
	Id         string
	Nombre     string
	Email      string
	Rol        string
	Biografia  string
	Web_url    string
	Web_titulo string
}

type ResumenPost struct {
	Id                string
	Fecha_publicacion string
	Alt_portada       string
	Titulo            string
	Resumen           string
	Autor_id          string
	Autor_nombre      string `db:"nombre"`
	Serie_id          string
	Serie_posicion    int
}

type FullPost struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string
	Alt_portada         string
	Titulo              string
	Resumen             string
	Contenido           string
	Autor_id            string
	Autor_nombre        string
	Autor_rol           string
	Autor_biografia     string
	Serie               sql.NullString
}

type Serie struct {
	Titulo    string
	Articulos []ResumenPost
}

type Tag struct {
	Id     string
	Nombre string
}
