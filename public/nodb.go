package public

import (
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
)

func SiguenosPage(w http.ResponseWriter, r *http.Request) {
	err := t.ExecuteTemplate(w, "siguenos.html", common.NoPageData{
		Meta: common.PageMeta{
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
	err := t.ExecuteTemplate(w, "licencias.html", common.NoPageData{
		Meta: common.PageMeta{
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
	err := t.ExecuteTemplate(w, "contacto.html", common.NoPageData{
		Meta: common.PageMeta{
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
