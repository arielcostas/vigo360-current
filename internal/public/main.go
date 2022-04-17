/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

import (
	"bytes"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
)

var db *sqlx.DB

func FullCanonica(path string) string {
	return os.Getenv("DOMAIN") + path
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	var output bytes.Buffer
	err := t.ExecuteTemplate(&output, "_404.html", NoPageData{
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

	w.Write(output.Bytes())
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	var output bytes.Buffer
	err := t.ExecuteTemplate(&output, "_500.html", NoPageData{
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
	w.Write(output.Bytes())
}

func InitRouter() *mux.Router {
	db = database.GetDB()

	router := mux.NewRouter().StrictSlash(true)

	router.Handle(`/post/{postid:[A-Za-z0-9\-\_|ñ]+}`, appHandler(PostPage)).Methods("GET")

	router.Handle(`/tags`, appHandler(listTags)).Methods(http.MethodGet)
	router.Handle(`/tags/{tagid:[0-9]+}/`, appHandler(viewTag)).Methods("GET")

	router.Handle(`/papers/{.*}`, http.RedirectHandler("/trabajos", http.StatusMovedPermanently)).Methods("GET")
	router.Handle(`/trabajos`, appHandler(listTrabajos)).Methods("GET")
	router.Handle(`/trabajos/{trabajoid:[A-Za-z0-9\-\_|ñ]+}`, appHandler(viewTrabajo)).Methods("GET")

	router.Handle(`/authors/{.*}`, http.RedirectHandler("/autores", http.StatusMovedPermanently)).Methods("GET")
	router.Handle(`/autores/{id:[A-Za-z0-9\-\_|ñ]+}`, appHandler(AutoresIdPage)).Methods("GET")
	router.Handle(`/autores`, appHandler(AutoresPage)).Methods("GET")

	router.Handle(`/siguenos`, appHandler(NoDbPage)).Methods("GET")
	router.Handle(`/legal`, appHandler(NoDbPage)).Methods("GET")
	router.Handle(`/contacto`, appHandler(NoDbPage)).Methods("GET")

	router.HandleFunc(`/sitemap.xml`, GenerateSitemap).Methods("GET")
	router.Handle(`/atom.xml`, appHandler(PostsAtomFeed)).Methods("GET")
	router.HandleFunc(`/trabajos/atom.xml`, TrabajosAtomFeed).Methods("GET")
	router.Handle(`/tags/{tagid:[0-9]+}/atom.xml`, appHandler(TagsAtomFeed)).Methods("GET")
	router.Handle(`/autores/{autorid}/atom.xml`, appHandler(AutorAtomFeed)).Methods("GET")

	router.Handle("/", appHandler(indexPage)).Methods("GET")

	router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	return router
}
