package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handleAdminEditWorkPage() http.HandlerFunc {
	type returnParams struct {
		Work    models.Trabajo
		Session models.Session
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		postId := mux.Vars(r)["id"]

		var trabajo models.Trabajo

		trabajo, err := s.store.trabajo.ObtenerPorId(postId, false)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Error("trabajo no encontrado: %s", err.Error())
				s.handleError(r, w, 400, messages.ErrorPaginaNoEncontrada)
			} else {
				log.Error("error recuperando trabajo: %s", err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
			}

			return
		}

		err = templates.Render(w, "admin-works-id.html", returnParams{
			Work:    trabajo,
			Session: sess,
		})
		if err != nil {
			log.Error("error mostrando página: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}
}
