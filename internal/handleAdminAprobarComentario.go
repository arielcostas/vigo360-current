
package internal

import (
	"fmt"
	"net/http"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/service"
)

func (s *Server) handleAdminAprobarComentario() http.HandlerFunc {
	var cs = service.NewComentarioService(s.store.comentario, s.store.publicacion)

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)
		cid := r.URL.Query().Get("cid") // cid = commentId = el comentario a rechazar

		err := cs.Aprobar(cid, sess.Autor_id)
		if err != nil {
			logger.Error("error aprobando comentario %s: %s", cid, err.Error())
			fmt.Fprintf(w, "Hubo un error aprobando el comentario")
		}

		w.Header().Add("Location", "/admin/comentarios")
		defer w.WriteHeader(303)
	}
}
