/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Server struct {
	Router *mux.Router
	store  *Container
}

func NewServer(c *Container) *Server {
	s := &Server{
		store: c,
	}
	s.Routes()
	return s
}

type ridContextKey string

func (s *Server) Routes() {
	if s.Router != nil {
		return
	}

	s.Router = mux.NewRouter().StrictSlash(true)
	s.Router.Use(middleware)

	s.Router.Handle("/admin/", http.RedirectHandler("/admin/login", http.StatusFound)).Methods(http.MethodGet)
	s.Router.HandleFunc("/admin/login", s.handleAdminLoginPage("")).Methods(http.MethodGet)
	s.Router.HandleFunc("/admin/login", s.handleAdminLoginAction()).Methods(http.MethodPost)
	s.Router.HandleFunc("/admin/logout", s.withAuth(s.handleAdminLogoutAction())).Methods(http.MethodGet)

	s.Router.HandleFunc("/admin/comentarios", s.withAuth(s.handleAdminListComentarios())).Methods(http.MethodGet)

	s.Router.HandleFunc("/admin/dashboard", s.withAuth(s.handleAdminDashboardPage())).Methods(http.MethodGet)
	s.Router.HandleFunc("/admin/post", s.withAuth(s.handleAdminListPost())).Methods(http.MethodGet)
	s.Router.HandleFunc("/admin/post", s.withAuth(s.handleAdminCreatePost())).Methods(http.MethodPost)

	s.Router.HandleFunc("/admin/post/{id}", s.withAuth(s.handleAdminEditPage())).Methods(http.MethodGet)
	s.Router.HandleFunc("/admin/post/{id}", s.withAuth(s.handleAdminEditAction())).Methods(http.MethodPost)
	s.Router.HandleFunc("/admin/post/{postid}/delete", s.withAuth(s.handleAdminDeletePost())).Methods(http.MethodGet)

	s.Router.HandleFunc("/admin/series", s.withAuth(s.handleAdminListSeries())).Methods(http.MethodGet)
	s.Router.HandleFunc("/admin/series", s.withAuth(s.handleAdminCreateSeries())).Methods(http.MethodPost)

	s.Router.HandleFunc("/admin/perfil", s.withAuth(s.handleAdminPerfilView())).Methods(http.MethodGet)
	s.Router.HandleFunc("/admin/perfil", s.withAuth(s.handleAdminPerfilEdit())).Methods(http.MethodPost)

	s.Router.HandleFunc("/admin/preview", s.withAuth(s.handleAdminPreviewPage())).Methods(http.MethodPost)

	s.Router.HandleFunc("/admin/async/fotosExtra", s.withAuth(s.handleAdminListarFotoExtra())).Methods(http.MethodGet)
	s.Router.HandleFunc("/admin/async/fotosExtra", s.withAuth(s.handleAdminCrearFotoExtra())).Methods(http.MethodPost)
	s.Router.HandleFunc("/admin/async/fotosExtra", s.withAuth(s.handleAdminDeleteFotoExtra())).Methods(http.MethodDelete)

	s.Router.HandleFunc(`/post/{postid}`, s.handlePublicPostPage()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/post/{postid}`, s.handlePublicEnviarComentario()).Methods(http.MethodPost)

	s.Router.HandleFunc(`/tags`, s.handlePublicListTags()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/tags/{tagid}/`, s.handlePublicTagPage()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/trabajos`, s.handlePublicListTrabajos()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/trabajos/{trabajoid}`, s.handlePublicTrabajoPage()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/autores/{id}`, s.handlePublicAutorPage()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/autores`, s.handlePublicListAutores()).Methods(http.MethodGet)

	s.Router.HandleFunc(`/legal`, s.handlePublicNodbPage()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/contacto`, s.handlePublicNodbPage()).Methods(http.MethodGet)

	s.Router.HandleFunc(`/atom.xml`, s.handlePublicIndexAtom()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/tags/{tagid}/atom.xml`, s.handlePublicTagsAtom()).Methods(http.MethodGet)
	s.Router.HandleFunc(`/autores/{autorid}/atom.xml`, s.handlePublicAutorAtom()).Methods(http.MethodGet)

	s.Router.HandleFunc(`/sitemap.xml`, s.handlePublicSitemap()).Methods(http.MethodGet)
	s.Router.HandleFunc("/buscar", s.handlePublicBusqueda()).Methods(http.MethodGet)

	var indexnowkeyurl = fmt.Sprintf("/%s.txt", os.Getenv("INDEXNOW_KEY"))
	s.Router.HandleFunc(indexnowkeyurl, s.handlePublicIndexnowKey()).Methods(http.MethodGet)

	s.Router.HandleFunc("/", s.handlePublicIndex()).Methods(http.MethodGet)
}
