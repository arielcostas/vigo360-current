package routes

import (
	"net/http"
)

type NoPageData struct{}

func SiguenosPage(w http.ResponseWriter, r *http.Request) {
	t.ExecuteTemplate(w, "siguenos.html", NoPageData{})
}
