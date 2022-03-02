package routes

import (
	"net/http"
)

func SiguenosPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "siguenos.html", NoPageData{
		Meta: PageMeta{Title: "Síguenos"},
	})
}

func LicenciasPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "licencias.html", NoPageData{
		Meta: PageMeta{Title: "Licencias"},
	})
}

func ContactoPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "contacto.html", NoPageData{
		Meta: PageMeta{Title: "Contacto"},
	})
}
