package public

import (
	"embed"
	"html/template"
	"net/http"
	"os"
	"strings"
	texttemplate "text/template"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

//go:embed html/*
var rawtemplates embed.FS

var t *template.Template
var tt *texttemplate.Template

func loadTemplates() {
	t = template.Must(template.ParseFS(rawtemplates, "html/*.html"))
	tt = texttemplate.Must(texttemplate.ParseFS(rawtemplates, "html/*.xml"))
}

func FullCanonica(path string) string {
	return os.Getenv("DOMAIN") + path
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	t.ExecuteTemplate(w, "_404.html", common.NoPageData{
		Meta: common.PageMeta{
			Titulo:      "Página no encontrada",
			Descripcion: "The requested resource could not be found in this server.",
			Canonica:    FullCanonica(r.URL.Path),
		},
	})
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	t.ExecuteTemplate(w, "_500.html", common.NoPageData{
		Meta: common.PageMeta{
			Titulo:      "Error del servidor",
			Descripcion: "There was a server error trying to load this page.",
			Canonica:    FullCanonica(r.URL.Path),
		},
	})
}

func AuthorsToAutores(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		strings.ReplaceAll(r.URL.String(), "/authors/", "/autores/"),
		http.StatusMovedPermanently)
}

func PapersToTrabajos(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		strings.ReplaceAll(r.URL.String(), "/papers/", "/trabajos/"),
		http.StatusMovedPermanently)
}

func InitRouter() *mux.Router {
	db = common.Database
	loadTemplates()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc(`/post/{postid:[A-Za-z0-9\-\_|ñ]+}`, PostPage).Methods("GET")

	router.HandleFunc(`/tags/{tagid:[0-9]+}`, TagsIdPage).Methods("GET")

	router.HandleFunc(`/papers/{.*}`, PapersToTrabajos).Methods("GET")
	router.HandleFunc(`/trabajos/{paperid:[A-Za-z0-9\-\_|ñ]+}`, TrabajoPage).Methods("GET")
	router.HandleFunc(`/trabajos`, TrabajoListPage).Methods("GET")

	router.HandleFunc(`/authors/{.*}`, AuthorsToAutores).Methods("GET")
	router.HandleFunc(`/autores/{id:[A-Za-z0-9\-\_|ñ]+}`, AutoresIdPage).Methods("GET")
	router.HandleFunc(`/autores`, AutoresPage).Methods("GET")

	router.HandleFunc(`/siguenos`, SiguenosPage).Methods("GET")
	router.HandleFunc(`/licencia`, LicenciasPage).Methods("GET")
	router.HandleFunc(`/contacto`, ContactoPage).Methods("GET")

	router.HandleFunc(`/sitemap.xml`, GenerateSitemap).Methods("GET")
	router.HandleFunc(`/atom.xml`, PostsAtomFeed).Methods("GET")

	router.HandleFunc("/", IndexPage).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	return router
}
