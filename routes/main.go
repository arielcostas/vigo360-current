package routes

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"os"
	"regexp"
	"strings"
	texttemplate "text/template"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type NoPageData struct {
	Meta PageMeta
}

type PageMeta struct {
	Titulo      string
	Descripcion string
	Canonica    string
	Miniatura   string
}

//go:embed html/*
var rawtemplates embed.FS

//go:embed includes
var includes embed.FS

var t *template.Template
var tt *texttemplate.Template

func loadTemplates() {
	t = template.Must(template.ParseFS(rawtemplates, "html/*.html"))
	tt = texttemplate.Must(texttemplate.ParseFS(rawtemplates, "html/*.xml"))
}

func FullCanonica(path string) string {
	return os.Getenv("DOMAIN") + path
}

func TestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func ValidatePassword(password string, hash string) bool {
	res := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if res != nil {
		println(res.Error())
		return false
	}
	return true
}

func includesHandler(w http.ResponseWriter, r *http.Request) {
	file := mux.Vars(r)["file"]
	ext := regexp.MustCompile(`\.[A-Za-z]+$`).FindString(file)
	bytes, err := includes.ReadFile("includes/" + file)
	if err != nil {
		// TODO error handling
		log.Fatalf(err.Error())
	}
	w.Header().Add("Content-Type", mime.TypeByExtension(ext))
	w.Header().Add("Cache-Control", "max-age=3600")
	w.Write(bytes)
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	t.ExecuteTemplate(w, "_404.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Página no encontrada",
			Descripcion: "The requested resource could not be found in this server.",
			Canonica:    FullCanonica(r.URL.Path),
		},
	})
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	t.ExecuteTemplate(w, "_500.html", NoPageData{
		Meta: PageMeta{
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
	InitDB()
	loadTemplates()

	router := mux.NewRouter().StrictSlash(true)
	router.Use(TestMiddleware)
	router.HandleFunc("/admin/login", AdminLogin).Methods("GET", "POST")

	router.HandleFunc(`/post/{postid:[A-Za-z0-9\-\_|ñ]+}`, PostPage).Methods("GET")

	router.HandleFunc(`/tags/{tagid:[0-9]+}`, TagsIdPage).Methods("GET")

	router.HandleFunc(`/papers/{.*}`, PapersToTrabajos).Methods("GET")
	router.HandleFunc(`/trabajos/{paperid:[A-Za-z0-9\-\_|ñ]+}`, TrabajoPage).Methods("GET")
	router.HandleFunc(`/trabajos`, TrabajoListPage).Methods("GET")

	router.HandleFunc(`/authors/{.*}`, AuthorsToAutores).Methods("GET")
	router.HandleFunc(`/autores/{id:[A-Za-z0-9\-\_|ñ]+}`, AutoresIdPage).Methods("GET")
	router.HandleFunc(`/autores`, AutoresPage).Methods("GET")

	router.HandleFunc(`/includes/{file:[\w|\.|\-|\_|ñ]+}`, includesHandler).Methods("GET")
	router.HandleFunc(`/siguenos`, SiguenosPage).Methods("GET")
	router.HandleFunc(`/licencia`, LicenciasPage).Methods("GET")
	router.HandleFunc(`/contacto`, ContactoPage).Methods("GET")

	router.HandleFunc(`/atom.xml`, PostsAtomFeed).Methods("GET")

	router.HandleFunc("/", IndexPage).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	return router
}
