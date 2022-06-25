/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"database/sql"
	"errors"
)

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

var ErrImageFormatError error = errors.New("invalid image MIME type")
