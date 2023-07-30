package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"vigo360.es/new/internal/messages"

	"github.com/gorilla/mux"
	"github.com/kataras/hcaptcha"
	"github.com/thanhpk/randstr"
)

func (s *Server) JsonifyRoutes(router *mux.Router, path string) *mux.Router {
	var newrouter = router
	newrouter.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var isJsonRoute = strings.HasPrefix(r.URL.Path, path)
			if isJsonRoute {
				w.Header().Add("Content-Type", "application/json")
			}
			h.ServeHTTP(w, r)
			if isJsonRoute {
				fmt.Fprintf(w, "\n")
			}
		})
	})
	return newrouter
}

func (s *Server) IdentifyRequests(router *mux.Router) *mux.Router {
	var newrouter = router
	newrouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var rid = randstr.String(15)
			fmt.Printf("<6>[%s] [%s] %s %s\n", r.Header.Get("X-Forwarded-For"), rid, r.Method, r.URL.Path)
			newContext := context.WithValue(r.Context(), ridContextKey("rid"), rid)
			r = r.WithContext(newContext)
			w.Header().Add("vigo360-rid", rid)
			next.ServeHTTP(w, r)
		})
	})
	return newrouter
}

func (s *Server) SetupApiRoutes(router *mux.Router) *mux.Router {
	var newrouter = router

	newrouter.HandleFunc("/api/v1/comentarios", s.withJsonAuth(s.handle_api_listar_comentarios)).Methods(http.MethodGet)

	return newrouter
}

func (s *Server) SetupWebRoutes(router *mux.Router) *mux.Router {
	var newrouter = router

	newrouter.Handle("/admin/", http.RedirectHandler("/admin/login", http.StatusFound)).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/login", s.handleAdminLoginPage("")).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/login", s.handleAdminLoginAction()).Methods(http.MethodPost)
	newrouter.HandleFunc("/admin/logout", s.withAuth(s.handleAdminLogoutAction())).Methods(http.MethodGet)

	newrouter.HandleFunc("/admin/comentarios", s.withAuth(s.handleAdminListComentarios())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/comentarios/aprobar", s.withAuth(s.handleAdminAprobarComentario())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/comentarios/rechazar", s.withAuth(s.handleAdminRechazarComentario())).Methods(http.MethodGet)

	newrouter.HandleFunc("/admin/dashboard", s.withAuth(s.handleAdminDashboardPage())).Methods(http.MethodGet)

	newrouter.HandleFunc("/admin/post", s.withAuth(s.handleAdminListPost())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/post", s.withAuth(s.handleAdminCreatePost())).Methods(http.MethodPost)
	newrouter.HandleFunc("/admin/post/{id}", s.withAuth(s.handleAdminEditPostPage())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/post/{id}", s.withAuth(s.handleAdminEditPostAction())).Methods(http.MethodPost)
	newrouter.HandleFunc("/admin/post/{postid}/delete", s.withAuth(s.handleAdminDeletePost())).Methods(http.MethodGet)

	newrouter.HandleFunc("/admin/works", s.withAuth(s.handleAdminListWorks())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/works", s.withAuth(s.handleAdminCreateWork())).Methods(http.MethodPost)
	newrouter.HandleFunc("/admin/works/{id}", s.withAuth(s.handleAdminEditWorkPage())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/works/{id}", s.withAuth(s.handleAdminEditWorkAction())).Methods(http.MethodPost)

	newrouter.HandleFunc("/admin/perfil", s.withAuth(s.handleAdminPerfilView())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/perfil", s.withAuth(s.handleAdminPerfilEdit())).Methods(http.MethodPost)

	newrouter.HandleFunc("/admin/preview", s.withAuth(s.handleAdminPreviewPage())).Methods(http.MethodPost)

	newrouter.HandleFunc("/admin/async/fotosExtra", s.withAuth(s.handleAdminListarFotoExtra())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/async/fotosExtra", s.withAuth(s.handleAdminCrearFotoExtra())).Methods(http.MethodPost)
	newrouter.HandleFunc("/admin/async/fotosExtra", s.withAuth(s.handleAdminDeleteFotoExtra())).Methods(http.MethodDelete)

	newrouter.HandleFunc(`/post/{postid}`, s.handlePublicPostPage()).Methods(http.MethodGet)

	var secretKey = os.Getenv("HCAPTCHA_SECRET")

	var cli = hcaptcha.New(secretKey)

	newrouter.HandleFunc(`/post/{postid}`, cli.HandlerFunc(s.handlePublicEnviarComentario())).Methods(http.MethodPost)

	newrouter.HandleFunc(`/tags`, s.handlePublicListTags()).Methods(http.MethodGet)
	newrouter.HandleFunc(`/tags/{tagid}/`, s.handlePublicTagPage()).Methods(http.MethodGet)
	newrouter.HandleFunc(`/trabajos`, s.handlePublicListTrabajos()).Methods(http.MethodGet)
	newrouter.HandleFunc(`/trabajos/{trabajoid}`, s.handlePublicTrabajoPage()).Methods(http.MethodGet)
	newrouter.HandleFunc(`/autores/{id}`, s.handlePublicAutorPage()).Methods(http.MethodGet)
	newrouter.HandleFunc(`/autores`, s.handlePublicListAutores()).Methods(http.MethodGet)

	newrouter.HandleFunc(`/legal`, s.handlePublicNodbPage()).Methods(http.MethodGet)
	newrouter.HandleFunc(`/contacto`, s.handlePublicNodbPage()).Methods(http.MethodGet)

	newrouter.HandleFunc(`/atom.xml`, s.handlePublicIndexAtom()).Methods(http.MethodGet)

	newrouter.HandleFunc(`/sitemap.xml`, s.handlePublicSitemap()).Methods(http.MethodGet)
	newrouter.HandleFunc("/buscar", s.handlePublicBusqueda()).Methods(http.MethodGet)

	var indexnowkeyurl = fmt.Sprintf("/%s.txt", os.Getenv("INDEXNOW_KEY"))
	newrouter.HandleFunc(indexnowkeyurl, s.handlePublicIndexnowKey()).Methods(http.MethodGet)

	newrouter.HandleFunc("/", s.handlePublicIndex()).Methods(http.MethodGet)

	newrouter.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.handleError(r, w, http.StatusNotFound, messages.ErrorPaginaNoEncontrada)
	})
	return newrouter
}
