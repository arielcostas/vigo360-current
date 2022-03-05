package public

import (
	"net/http"
)

func SiguenosPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "siguenos.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Síguenos",
			Descripcion: "Información sobre cómo seguir a Vigo360, y enterarse de sus últimas publicaciones y novedades.",
			Canonica:    FullCanonica("/siguenos"),
		},
	})
}

func LicenciasPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "licencias.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Licencias",
			Descripcion: "Información legal relativa a Vigo360, desde licencias de uso libre hasta la política de privacidad.",
			Canonica:    FullCanonica("/licencia"),
		},
	})
}

func ContactoPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "contacto.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Contacto",
			Descripcion: "Si necesitases contactar con Vigo360, aquí encontrarás cómo hacerlo.",
			Canonica:    FullCanonica("/contacto"),
		},
	})
}
