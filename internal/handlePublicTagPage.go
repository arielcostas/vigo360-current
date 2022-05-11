package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/model"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handlePublicTagPage() http.HandlerFunc {
	type response struct {
		Tag   model.Tag
		Posts model.Publicaciones
		Meta  PageMeta
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		req_tagid := mux.Vars(r)["tagid"]

		tag, err := s.store.tag.Obtener(req_tagid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró la tag: %s", err.Error())
				s.handleError(w, 404, messages.ErrorPaginaNoEncontrada)
			} else {
				logger.Error("error recuperando la tag: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}
			return
		}

		publicaciones, err := s.store.publicacion.ListarPorTag(req_tagid)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando publicaciones para tag: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		err = templates.Render(w, "tags-id.html", response{
			Tag:   tag,
			Posts: publicaciones.FiltrarPublicas(),
			Meta: PageMeta{
				Titulo:      tag.Nombre,
				Keywords:    tag.Nombre,
				Descripcion: "Publicaciones en Vigo360 sobre " + tag.Nombre,
				Canonica:    fullCanonica("/tags/" + req_tagid),
			},
		})

		if err != nil {
			logger.Error("error generando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
