/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

func SiguenosPage(w http.ResponseWriter, r *http.Request) {
	err := t.ExecuteTemplate(w, "siguenos.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Síguenos",
			Descripcion: "Información sobre cómo seguir a Vigo360, y enterarse de sus últimas publicaciones y novedades.",
			Canonica:    FullCanonica("/siguenos"),
		},
	})
	if err != nil {
		logger.Error("[siguenos] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
}

func LicenciasPage(w http.ResponseWriter, r *http.Request) {
	err := t.ExecuteTemplate(w, "licencias.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Licencias",
			Descripcion: "Información legal relativa a Vigo360, desde licencias de uso libre hasta la política de privacidad.",
			Canonica:    FullCanonica("/licencia"),
		},
	})
	if err != nil {
		logger.Error("[licencias] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
}

func ContactoPage(w http.ResponseWriter, r *http.Request) {
	err := t.ExecuteTemplate(w, "contacto.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Contacto",
			Descripcion: "Si necesitases contactar con Vigo360, aquí encontrarás cómo hacerlo.",
			Canonica:    FullCanonica("/contacto"),
		},
	})
	if err != nil {
		logger.Error("[contacto] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
}
