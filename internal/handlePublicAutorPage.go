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

func (s *Server) handlePublicAutorPage() http.HandlerFunc {
	type Response struct {
		Autor    models.Autor
		Posts    models.Publicaciones
		Trabajos models.Trabajos
		Meta     PageMeta
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var req_autor = mux.Vars(r)["id"]

		var autor models.Autor
		autor, err := s.store.autor.Obtener(req_autor)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró autor coon id "+req_autor, err.Error())
				s.handleError(r, w, 404, "No se ha encontrado el autor solicitado")
			} else {
				logger.Error("error inesperado recuperando datos: %s", err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
			}
			return
		}

		publicaciones, err := s.store.publicacion.ListarPorAutor(autor.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		trabajos, err := s.store.trabajo.ListarPorAutor(autor.Id)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		err = templates.Render(w, "autores-id.html", Response{
			Autor:    autor,
			Posts:    publicaciones.FiltrarPublicas(),
			Trabajos: trabajos,
			Meta: PageMeta{
				Titulo:      autor.Nombre,
				Descripcion: autor.Biografia,
				Canonica:    fullCanonica("/autores/" + autor.Id),
			},
		})

		if err != nil {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
		}
	}
}
