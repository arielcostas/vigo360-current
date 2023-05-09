package internal

import (
	"net/http"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
)

func (s *Server) handleAdminDeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)
		if !sess.Permisos["publicaciones_delete"] {
			logger.Error("sin permiso para eliminar la publicación")
			s.handleError(w, 403, messages.ErrorSinPermiso)
			return
		}

		// TODO: Convertir esto en procedimiento
		var postid = mux.Vars(r)["postid"]
		tx, err := database.GetDB().Begin()
		if err != nil {
			logger.Error("error iniciando transacción: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}
		_, err = tx.Exec("DELETE FROM publicaciones_tags WHERE publicacion_id=?", postid)
		if err != nil {
			logger.Error("error eliminando tags: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}
		_, err = tx.Exec("DELETE FROM publicaciones WHERE id=?", postid)
		if err != nil {
			logger.Error("error eliminando publicación: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}
		err = tx.Commit()
		if err != nil {
			logger.Error("error ejecutando transacción: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		w.Header().Add("Location", "/admin/post")
		defer w.WriteHeader(307)
	}
}
