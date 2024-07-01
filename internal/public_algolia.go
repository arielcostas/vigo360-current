package internal

import (
	"encoding/json"
	"net/http"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handlePublicIndexAlgolia() http.HandlerFunc {
	type Post struct {
		Id                  string `json:"id"`
		Alt_portada         string `json:"alt_portada"`
		Titulo              string `json:"titulo"`
		Resumen             string `json:"resumen"`
		Contenido           string `json:"contenido"`
		Fecha_publicacion   string `json:"fecha_publicacion"`
		Fecha_actualizacion string `json:"fecha_actualizacion"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var result []Post = make([]Post, 0)

		db := database.GetDB()
		rows, err := db.Query("SELECT id, alt_portada, titulo, resumen, contenido, fecha_publicacion, fecha_actualizacion FROM publicaciones")

		if err != nil {
			log.Error("error leyendo adjuntos: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorDatos)
			return
		}

		for rows.Next() {
			var post Post

			err = rows.Scan(&post.Id, &post.Alt_portada, &post.Titulo, &post.Resumen, &post.Contenido, &post.Fecha_publicacion, &post.Fecha_actualizacion)
			if err != nil {
				log.Error("error escaneando adjuntos: %s", err.Error())
				s.handleJsonError(r, w, 500, messages.ErrorDatos)
				return
			}
			result = append(result, post)
		}

		resbytes, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			log.Error("error escribiendo json de respuesta: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorRender)
			return
		}
		w.Write(resbytes)
	}
}
