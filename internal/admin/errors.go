/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import "errors"

var ErrInvalidFormInput = errors.New("provided data is not valid")
var ErrExpiredSession = errors.New("session was older than 6 hours and was revoked automatically")
var ErrInvalidSession = errors.New("session token is not valid")
var ErrUnablePermissions = errors.New("unable to get permissions for session")
var ErrLoginRequired = errors.New("login is required to load this page")
var ErrBadParam = errors.New("param not provided or not valid")

var LoginRequiredAppError = &appError{ErrLoginRequired, "login token not provided or not valid", "Es necesario iniciar sesión para acceder a esta página", 401}

// Creates an error related with template rendering
func newTemplateRenderingAppError(err error) *appError {
	return &appError{Error: err, Message: "error rendering template",
		Response: "Hubo un error intentando mostrar la página.", Status: 500}
}

// Creates an error related with fetching from database
func newDatabaseReadAppError(err error, datatype string) *appError {
	if len(datatype) == 0 {
		datatype = " "
	}
	return &appError{Error: err, Message: "error fetching " + datatype + "from database",
		Response: "Error leyendo datos", Status: 500}
}
