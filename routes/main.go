package routes

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"mime"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

//go:embed html/*
var rawtemplates embed.FS

//go:embed includes
var includes embed.FS

var t *template.Template

func loadTemplates() {
	t = template.Must(template.ParseFS(rawtemplates, "html/*.html"))
	t.Funcs(template.FuncMap{
		"formatDate": func(date []uint8) string {
			println(date)
			return string(date)
		},
	})
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
		log.Fatalf(err.Error())
	}
	w.Header().Add("Content-Type", mime.TypeByExtension(ext))
	w.Header().Add("Cache-Control", "max-age=3600")
	w.Write(bytes)
}

func InitRouter() *mux.Router {
	InitDB()
	loadTemplates()

	router := mux.NewRouter().StrictSlash(true)
	router.Use(TestMiddleware)
	router.HandleFunc("/admin/login", AdminLogin).Methods("GET", "POST")
	router.HandleFunc(`/post/{postid:[A-Za-z0-9\-\_|ñ]+}`, PostPage).Methods("GET")
	router.HandleFunc(`/tags/{tagid:[0-9]+}`, TagsIdPage).Methods("GET")
	router.HandleFunc(`/includes/{file:[\w|\.|\-|\_|ñ]+}`, includesHandler).Methods("GET")
	router.HandleFunc("/", IndexPage).Methods("GET")
	return router
}
