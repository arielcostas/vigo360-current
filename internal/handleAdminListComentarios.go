package internal

import (
	"net/http"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handleAdminListComentarios() http.HandlerFunc {
	type Response struct {
		Comentarios []models.Comentario
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		//sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)
		comentarios, err := s.store.comentario.ListarPorEstado(models.EstadoPendiente)
		if err != nil {
			log.Error("Error recuperando comentarios: " + err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		err = templates.Render(w, "admin-comentarios.html", Response{
			Comentarios: comentarios,
		})

		if err != nil {
			log.Error("error recuperando el autor: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
