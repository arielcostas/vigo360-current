package internal

import (
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/service"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicPostPage() http.HandlerFunc {
	type Response struct {
		Post            models.Publicacion
		LoggedIn        bool
		Comentarios     []service.ComentarioTree
		Recommendations []Sugerencia
		Meta            PageMeta
		HcaptchaClient  string
	}

	var cs = service.NewComentarioService(s.store.comentario, s.store.publicacion)

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var sc, err = r.Cookie("sess")
		var loggedIn = false
		if err == nil {
			_, err := s.getSession(sc.Value)

			if err == nil { // User is logged in
				loggedIn = true
			}
		}

		req_post_id := mux.Vars(r)["postid"]

		var post models.Publicacion
		if np, err := s.store.publicacion.ObtenerPorId(req_post_id, true); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró la publicación: %s", err.Error())
				s.handleError(w, 404, messages.ErrorPaginaNoEncontrada)
			} else {
				logger.Error("error recuperando la publicación: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}
			return
		} else {
			post = np
		}

		if post.Serie.Id != "" {
			var err error
			post.Serie, err = s.store.serie.Obtener(post.Serie.Id)
			if err != nil {
				logger.Error("error recuperando serie de la publicación: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}
		}

		var recommendations []Sugerencia
		if nr, err := generateSuggestions(post.Id, s.store.publicacion); err != nil {
			logger.Error("error recuperando sugerencias: %s", err.Error())
			recommendations = make([]Sugerencia, 0)
		} else {
			recommendations = nr
		}

		var keywords = ""
		for _, t := range post.Tags {
			keywords += t.Nombre + ","
		}

		post.Serie.Publicaciones = post.Serie.Publicaciones.FiltrarPublicas()
		var ct []service.ComentarioTree

		if nct, err := cs.ListarPublicos(post.Id); err != nil {
			logger.Error("error recuperando comentarios para %s: %s", post.Id, err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
		} else {
			ct = nct
		}

		err = templates.Render(w, "post-id.html", Response{
			Post:            post,
			LoggedIn:        loggedIn,
			Recommendations: recommendations,
			Comentarios:     ct,
			HcaptchaClient:  os.Getenv("HCAPTCHA_SITEKEY"),
			Meta: PageMeta{
				Titulo:      post.Titulo,
				Descripcion: post.Resumen,
				Keywords:    keywords,
				Canonica:    fullCanonica("/post/" + post.Id),
				Miniatura:   fullCanonica("/static/thumb/" + post.Id + ".jpg"),
			},
		})
		if err != nil {
			logger.Error("error mostrando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
