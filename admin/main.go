/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"embed"
	"html/template"
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/common"
	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

//go:embed html/*
var rawtemplates embed.FS

var t = func() *template.Template {
	t := template.New("")

	functions := template.FuncMap{
		"safeHTML": func(text string) template.HTML {
			return template.HTML(text)
		},
		"sum": func(a int, b int) int {
			return a + b
		},
	}

	entries, _ := rawtemplates.ReadDir("html")
	for _, de := range entries {
		filename := de.Name()
		contents, _ := rawtemplates.ReadFile("html/" + filename)

		_, err := t.New(filename).Funcs(functions).Parse(string(contents))
		if err != nil {
			logger.Critical("[public-main] error parsing template: %s", err.Error())
		}
	}

	return t
}()

var db *sqlx.DB

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

func redirectToDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/admin/login")
	w.WriteHeader(302)
}

func InitRouter() *mux.Router {
	db = common.Database

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/admin/", redirectToDashboard).Methods("GET")
	router.HandleFunc("/admin/login", LoginPage).Methods("GET")
	router.HandleFunc("/admin/login", LoginAction).Methods("POST")
	router.HandleFunc("/admin/logout", LogoutPage).Methods("GET")

	router.HandleFunc("/admin/dashboard", DashboardPage).Methods("GET")
	router.HandleFunc("/admin/post", PostListPage).Methods("GET")
	router.Handle("/admin/post", appHandler(createPost)).Methods("POST")

	router.HandleFunc("/admin/post/{id}", EditPostPage).Methods("GET")
	router.HandleFunc("/admin/post/{id}", EditPostAction).Methods("POST")

	router.HandleFunc("/admin/series", ListSeries).Methods("GET")
	router.HandleFunc("/admin/series", CreateSeries).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	return router
}
