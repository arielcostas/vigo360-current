package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handleAdminEditPage() http.HandlerFunc {
	type tag struct {
		models.Tag
		Seleccionada bool
	}

	type returnParams struct {
		Post    models.Publicacion
		Tags    []tag
		Session models.Session
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		postId := mux.Vars(r)["id"]

		db := database.GetDB()

		var publicacion models.Publicacion

		publicacion, err := s.store.publicacion.ObtenerPorId(postId, false)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				log.Error("publicacion no encontrada: %s", err.Error())
				s.handleError(r, w, 400, messages.ErrorPaginaNoEncontrada)
			} else {
				log.Error("error recuperando publicacion: %s", err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
			}

			return
		}

		var tags []tag

		//goland:noinspection SqlConstantExpression
		err = db.Select(&tags, `SELECT id, nombre, (SELECT tag_id FROM publicaciones_tags pt WHERE pt.publicacion_id = ? AND pt.tag_id = id) IS NOT NULL as seleccionada FROM tags ORDER BY nombre`, postId)
		if err != nil {
			log.Error("error recuperando tags: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
		}

		err = templates.Render(w, "admin-post-id.html", returnParams{
			Post:    publicacion,
			Tags:    tags,
			Session: sess,
		})
		if err != nil {
			log.Error("error mostrando página: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}
}
