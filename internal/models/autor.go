// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package models

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
