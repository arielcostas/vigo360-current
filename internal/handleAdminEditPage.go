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
	type returnParams struct {
		Post    PostEditar
		Series  []Serie
		Tags    []Tag
		Session models.Session
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		post_id := mux.Vars(r)["id"]
		post := PostEditar{}
		db := database.GetDB()

		err := db.QueryRowx(`SELECT id, titulo, resumen, contenido, alt_portada, (fecha_publicacion is not null && fecha_publicacion < NOW()) as publicado, serie_id, serie_posicion FROM publicaciones WHERE id = ?;`, post_id).StructScan(&post)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("publicacion no encontrada: %s", err.Error())
				s.handleError(w, 400, messages.ErrorPaginaNoEncontrada)
			} else {
				logger.Error("error recuperando publicacion: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}

			return
		}

		series := []Serie{}
		err = db.Select(&series, `SELECT * FROM series;`)

		if err != nil {
			logger.Error("error recuperando series: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		tags := []Tag{}
		err = db.Select(&tags, `SELECT id, nombre, (SELECT tag_id FROM publicaciones_tags pt WHERE pt.publicacion_id = ? AND pt.tag_id = id) IS NOT NULL as seleccionada FROM tags ORDER BY nombre ASC`, post_id)
		if err != nil {
			logger.Error("error recuperando tags: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
		}

		err = templates.Render(w, "admin-post-id.html", returnParams{
			Post:    post,
			Series:  series,
			Tags:    tags,
			Session: sess,
		})
		if err != nil {
			logger.Error("error mostrando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
