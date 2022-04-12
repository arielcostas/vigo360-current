/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

type Web struct {
	Url    string
	Titulo string
}

type Autor struct {
	Id        string
	Nombre    string
	Email     string
	Rol       string
	Biografia string
	Web       Web

	Publicaciones Publicaciones
}
