package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"vigo360.es/new/internal/database"
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
				w.Header().Add("Content-Type", "application/json; charset=utf-8")
			}
			h.ServeHTTP(w, r)
			if isJsonRoute {
				fmt.Fprintf(w, "\n")
			}
		})
	})
	return newrouter
}

func (s *Server) IdentifySessions(router *mux.Router) *mux.Router {
	var newrouter = router
	newrouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var sid string
			var sidCookie, err = r.Cookie("sid")
			if sidCookie == nil || err != nil {
				sid = randstr.String(15)
				println("sid: " + sid)
				http.SetCookie(w, &http.Cookie{
					Name:     "sid",
					Value:    sid,
					Path:     "/",
					HttpOnly: true,
				})
			} else {
				println("Reusing sid: " + sidCookie.Value)
				sid = sidCookie.Value
			}

			newContext := context.WithValue(r.Context(), ridContextKey("sid"), sid)
			r = r.WithContext(newContext)
			next.ServeHTTP(w, r)
		})
	})
	return newrouter
}

func (s *Server) IdentifyRequests(router *mux.Router) *mux.Router {
	var newrouter = router
	newrouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var rid = randstr.String(15)
			newContext := context.WithValue(r.Context(), ridContextKey("rid"), rid)
			r = r.WithContext(newContext)
			w.Header().Add("vigo360-rid", rid)
			next.ServeHTTP(w, r)
		})
	})
	return newrouter
}

func (s *Server) LogRequests(router *mux.Router) *mux.Router {
	var newrouter = router
	newrouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var startTime = time.Now()

			next.ServeHTTP(w, r)

			var rid = r.Context().Value(ridContextKey("rid")).(string)
			var sid = r.Context().Value(ridContextKey("sid")).(string)
			var duration = time.Since(startTime).Milliseconds()
			var ip = r.Header.Get("X-Forwarded-For")
			var method = r.Method
			var path = r.URL.Path
			var ua = r.Header.Get("User-Agent")

			var db = database.GetDB()

			var query = `INSERT INTO log (rid, sid, time, ip ,url, method, time_taken_ms, user_agent) VALUES (?,?, ?, ?, ?, ?, ?, ?)`

			_, err := db.Exec(query, rid, sid, startTime, ip, path, method, duration, ua)
			if err != nil {
				fmt.Println(err)
			}
		})
	})
	return newrouter
}

func (s *Server) SetupWebRoutes(router *mux.Router) *mux.Router {
	var newrouter = router

	newrouter.Handle("/admin/", http.RedirectHandler("/admin/login", http.StatusFound)).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/login", s.handle_login_page("")).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/login", s.handle_login_action()).Methods(http.MethodPost)
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

	newrouter.HandleFunc("/admin/async/attachments", s.withAuth(s.adminApiAttachmentList())).Methods(http.MethodGet)
	newrouter.HandleFunc("/admin/async/attachments", s.withAuth(s.adminApiAttachmentCreate())).Methods(http.MethodPost)
	newrouter.HandleFunc("/admin/async/attachments", s.withAuth(s.adminApiAttachmentDelete())).Methods(http.MethodDelete)

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
