/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"net/http"
	"os"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func FullCanonica(path string) string {
	return os.Getenv("DOMAIN") + path
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	err := t.ExecuteTemplate(w, "_404.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Página no encontrada",
			Descripcion: "The requested resource could not be found in this server.",
			Canonica:    FullCanonica(r.URL.Path),
		},
	})

	if err != nil {
		logger.Error("[main] error rendering 404 page: %s", err.Error())
		w.Write([]byte("La página solicitada no fue encontrada. Adicionalmente, no fue posible mostrar la página de error correspondiente."))
		return
	}
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	err := t.ExecuteTemplate(w, "_500.html", NoPageData{
		Meta: PageMeta{
			Titulo:      "Error del servidor",
			Descripcion: "There was a server error trying to load this page.",
			Canonica:    FullCanonica(r.URL.Path),
		},
	})
	if err != nil {
		logger.Error("[main] error rendering 500 page (ironic): %s", err.Error())
		w.Write([]byte("Error interno del servidor. Adicionalmente, la página de error no puede ser mostrada."))
		return
	}
}

func AuthorsToAutores(w http.ResponseWriter, r *http.Request) *appError {
	w.WriteHeader(301)
	w.Header().Add("Location", "/autores/")

	return nil
}

func PapersToTrabajos(w http.ResponseWriter, r *http.Request) *appError {
	w.WriteHeader(301)
	w.Header().Add("Location", "/trabajos/")

	return nil
}

func InitRouter(database *sqlx.DB) *mux.Router {
	db = database

	router := mux.NewRouter().StrictSlash(true)

	router.Handle(`/post/{postid:[A-Za-z0-9\-\_|ñ]+}`, appHandler(PostPage)).Methods("GET")

	router.Handle(`/tags`, appHandler(listTags)).Methods(http.MethodGet)
	router.Handle(`/tags/{tagid:[0-9]+}/`, appHandler(viewTag)).Methods("GET")

	router.Handle(`/papers/{.*}`, appHandler(PapersToTrabajos)).Methods("GET")
	router.Handle(`/trabajos`, appHandler(listTrabajos)).Methods("GET")
	router.Handle(`/trabajos/{trabajoid:[A-Za-z0-9\-\_|ñ]+}`, appHandler(viewTrabajo)).Methods("GET")

	router.Handle(`/authors/{.*}`, appHandler(AuthorsToAutores)).Methods("GET")
	router.Handle(`/autores/{id:[A-Za-z0-9\-\_|ñ]+}`, appHandler(AutoresIdPage)).Methods("GET")
	router.Handle(`/autores`, appHandler(AutoresPage)).Methods("GET")

	router.Handle(`/siguenos`, appHandler(SiguenosPage)).Methods("GET")
	router.Handle(`/licencia`, appHandler(LicenciasPage)).Methods("GET")
	router.Handle(`/contacto`, appHandler(ContactoPage)).Methods("GET")

	router.HandleFunc(`/sitemap.xml`, GenerateSitemap).Methods("GET")
	router.HandleFunc(`/atom.xml`, PostsAtomFeed).Methods("GET")
	router.HandleFunc(`/trabajos/atom.xml`, TrabajosAtomFeed).Methods("GET")
	router.HandleFunc(`/tags/{tagid:[0-9]+}/atom.xml`, TagsAtomFeed).Methods("GET")
	router.HandleFunc(`/autores/{autorid}/atom.xml`, AutorAtomFeed).Methods("GET")

	router.Handle("/", appHandler(indexPage)).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	return router
}
