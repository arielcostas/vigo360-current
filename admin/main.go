package admin

import (
	"embed"
	"html/template"
	"log"
	"net/http"

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

func verifyLogin(w http.ResponseWriter, r *http.Request) SesionRow {
	// TODO error handling
	cookie, err := r.Cookie("sess")
	if err == http.ErrNoCookie && r.URL.Path != "/admin/login" {
		http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
		return SesionRow{}
	}

	user := SesionRow{}

	err = db.QueryRowx("SELECT sessid, id, nombre, rol FROM sesiones LEFT JOIN autores ON sesiones.autor_id = autores.id WHERE sessid = ? AND revocada = false;", cookie.Value).StructScan(&user)

	if err != nil && r.URL.Path != "/admin/login" {
		log.Println("error in login verification: " + err.Error())
		http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
		return SesionRow{}
	}

	// Logged in successfully, no sense to log in again
	if r.URL.Path == "/admin/login" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
	}

	return user
}

func InitRouter() *mux.Router {
	loadTemplates()
	db = common.Database

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/admin/login", LoginPage).Methods("GET")
	router.HandleFunc("/admin/login", LoginAction).Methods("POST")

	router.HandleFunc("/admin/dashboard", DashboardPage).Methods("GET", "POST")

	return router
}
