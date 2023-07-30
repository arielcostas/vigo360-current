package internal

import (
	"net/http"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/service"
)

func (s *Server) handlePublicEnviarComentario() http.HandlerFunc {
	var cs = service.NewComentarioService(s.store.comentario, s.store.publicacion)

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))

		r.ParseForm()
		var publicacion_id = mux.Vars(r)["postid"]
		var nombre = r.Form.Get("nombre")
		var contenido = r.Form.Get("contenido")
		var padre = r.Form.Get("padre")

		var es_autor = false
		var autor_original = false

		publicacion, err := s.store.publicacion.ObtenerPorId(publicacion_id, true)
		if err != nil {
			logger.Error("publicación %s no comentable: %s", publicacion_id, err.Error())
		}

		sc, err := r.Cookie("sess")
		if err == nil {
			sess, err := s.getSession(sc.Value)

			if err == nil { // User is logged in
				if publicacion.Autor.Id == sess.Autor_id {
					autor_original = true
				}
				es_autor = true
				nombre = sess.Autor_nombre
			}
		}

		var nc models.Comentario

		if padre == "" {
			nc, err = cs.AgregarComentario(publicacion_id, nombre, contenido, es_autor, autor_original)
		} else {
			nc, err = cs.AgregarRespuesta(publicacion_id, nombre, contenido, padre, es_autor, autor_original)
		}

		if err != nil {
			logger.Error("error guardando comentario: %s", err.Error())
			s.handleError(r, w, 400, messages.ErrorDatos)
			return
		}

		logger.Information("guardado comentario con ID %s", nc.Id)

		w.Header().Add("Location", r.URL.Path)
		defer w.WriteHeader(http.StatusSeeOther)
	}
}
