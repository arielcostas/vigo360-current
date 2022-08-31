// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package messages

type ErrorMessage string

var ErrorDatos ErrorMessage = "Hubo un error recuperando los datos."
var ErrorNoResultados ErrorMessage = "No se ha encontrado ningún resultado."
var ErrorPaginaNoEncontrada ErrorMessage = "No se ha encontrado ningún resultado."
var ErrorRender ErrorMessage = "Hubo un error mostrando la página solicitada. Inténtelo de nuevo más tarde."
var ErrorFormulario ErrorMessage = "Hubo un error recuperando los datos enviados."
var ErrorValidacion ErrorMessage = "Alguno de los datos del formulario no es válido"
var ErrorSinPermiso ErrorMessage = "No tienes permiso para ver esta página o realizar esta opción."
var ErrorSinAutenticar ErrorMessage = "Debes autenticarte antes de acceder a esta página"
