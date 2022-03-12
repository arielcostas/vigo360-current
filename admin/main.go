package admin

import (
	"database/sql"
	"embed"
	"errors"
	"html/template"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

//go:embed html/*
var rawtemplates embed.FS

var t *template.Template
var db *sqlx.DB

func loadTemplates() {
	var err error
	t, err = template.ParseFS(rawtemplates, "html/*.html")
	if err != nil {
		logger.Critical("error loading admin templates: %s", err.Error())
	}
}

func gotoLogin(w http.ResponseWriter, r *http.Request) Sesion {
	http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
	return Sesion{}
}

// TODO Refactor this
func verifyLogin(w http.ResponseWriter, r *http.Request) Sesion {
	cookie, err := r.Cookie("sess")
	if errors.Is(err, http.ErrNoCookie) && r.URL.Path != "/admin/login" {
		logger.Notice("unauthenticated user tried accessing auth-requiring page %s", r.URL.Path)
		return gotoLogin(w, r)
	} else if err != nil && r.URL.Path != "/admin/login" {
		logger.Error("error getting session cookie: %s", err.Error())
		return gotoLogin(w, r)
	} else if err != nil {
		logger.Error("error getting session cookie: %s", err.Error())
		return Sesion{}
	}

	user := Sesion{}

	err = db.QueryRowx("SELECT id, nombre, rol FROM sesiones LEFT JOIN autores ON sesiones.autor_id = autores.id WHERE sessid = ? AND revocada = false;", cookie.Value).StructScan(&user)

	if errors.Is(err, sql.ErrNoRows) && r.URL.Path != "/admin/login" {
		logger.Warning("error in login verification: %s", err.Error())
		return gotoLogin(w, r)
	} else if err != nil {
		logger.Error("unexpected error fetching session from database: %s", err.Error())
		return gotoLogin(w, r)
	}

	// Logged in successfully, no sense to log in again
	if r.URL.Path == "/admin/login" {
		http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
	}

	// It's not the login page and the user is logged in
	return user
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	w.WriteHeader(404)
	t.ExecuteTemplate(w, "_404.html", struct{}{})
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	w.WriteHeader(500)
	t.ExecuteTemplate(w, "_500.html", struct{}{})
}

func InitRouter() *mux.Router {
	loadTemplates()
	db = common.Database

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/admin/login", LoginPage).Methods("GET")
	router.HandleFunc("/admin/login", LoginAction).Methods("POST")

	router.HandleFunc("/admin/dashboard", DashboardPage).Methods("GET")
	router.HandleFunc("/admin/post", PostListPage).Methods("GET")
	router.HandleFunc("/admin/post", CreatePostAction).Methods("POST")

	router.HandleFunc("/admin/post/{id}", EditPostPage).Methods("GET")
	router.HandleFunc("/admin/post/{id}", EditPostAction).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	return router
}
