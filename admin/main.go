package admin

import (
	"embed"
	"html/template"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

//go:embed html/*
var rawtemplates embed.FS

var t *template.Template
var db *sqlx.DB

func loadTemplates() {
	t = template.Must(template.ParseFS(rawtemplates, "html/*.html"))
}

func InitRouter() *mux.Router {
	loadTemplates()
	db = common.Database

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/admin/login", LoginPage).Methods("GET")
	router.HandleFunc("/admin/login", LoginAction).Methods("POST")

	return router
}
