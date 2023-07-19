package internal

import (
	"fmt"
	"net/http"
	"time"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handleAdminPreviewPage() http.HandlerFunc {
	type response struct {
		Post  models.Publicacion
		Ahora string
	}

	//var aBase64 = func(b []byte) string {
	//	return base64.StdEncoding.EncodeToString(b)
	//}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))

		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		autor, err := s.store.autor.Obtener(sess.Autor_id)
		if err != nil {
			logger.Error("error obteniendo datos del autor: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		post := models.Publicacion{
			Id:                  r.FormValue("post-id"),
			Titulo:              r.FormValue("art-titulo"),
			Resumen:             r.FormValue("art-resumen"),
			Contenido:           r.FormValue("art-contenido"),
			Alt_portada:         r.FormValue("alt-portada"),
			Fecha_actualizacion: time.Now().Format("2006-01-02 15:04:05"),
			Fecha_publicacion:   time.Now().Format("2006-01-02 15:04:05"),
			Autor:               autor,
		}

		err = templates.Render(w, "admin-preview.html", response{
			Post:  post,
			Ahora: time.Now().Format("02/01 15:04:03 -07:00"),
		})

		if err != nil {
			// TODO: Remplazar esto
			fmt.Printf("%s\n", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
			return
		}
	}
}
