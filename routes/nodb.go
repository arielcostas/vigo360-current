package routes

import (
	"net/http"
)

type NoPageData struct{}

func SiguenosPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "siguenos.html", NoPageData{})
}

func LicenciasPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "licencias.html", NoPageData{})
}

func ContactoPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "contacto.html", NoPageData{})
}
