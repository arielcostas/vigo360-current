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

func (s *Server) handlePublicTrabajoPage() http.HandlerFunc {
	// TODO: Refactor esto
	type Adjunto struct {
		Nombre_archivo string
		Titulo         string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		trabajoid := mux.Vars(r)["trabajoid"]

		var trabajo models.Trabajo

		if nt, err := s.store.trabajo.ObtenerPorId(trabajoid, true); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("trabajo no encontrado: %s", err.Error())
				s.handleError(r, w, 404, messages.ErrorPaginaNoEncontrada)
			} else {
				logger.Error("error recuperando trabajo: %s", err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
			}
			return
		} else {
			trabajo = nt
		}

		var adjuntos = make([]Adjunto, 0)
		// TODO: Refactor esto
		err := database.GetDB().Select(&adjuntos, "SELECT nombre_archivo, titulo FROM adjuntos WHERE trabajo_id = ?;", trabajo.Id)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando adjuntos: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
		}

		err = templates.Render(w, "trabajos-id.html", struct {
			Trabajo  models.Trabajo
			Adjuntos []Adjunto
			Meta     PageMeta
		}{
			Trabajo:  trabajo,
			Adjuntos: adjuntos,
			Meta: PageMeta{
				Titulo:      trabajo.Titulo,
				Descripcion: trabajo.Resumen,
				Canonica:    fullCanonica("/trabajos/" + trabajo.Id),
				Miniatura:   fullCanonica("/static/thumb/" + trabajo.Id + ".jpg"),
				BaseUrl:     baseUrl(),
			},
		})

		if err != nil {
			logger.Error("error mostrando página: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}
}
