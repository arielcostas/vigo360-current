/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

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
