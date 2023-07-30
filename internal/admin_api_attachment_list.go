package internal

import (
	"encoding/json"
	"net/http"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) adminApiAttachmentList() http.HandlerFunc {
	type Attachment struct {
		Id       int    `json:"id"`
		Title    string `json:"title"`
		Filename string `json:"filename"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))

		trabajoId := r.URL.Query().Get("trabajo")
		if trabajoId == "" {
			log.Error("no se especificó un trabajo")
			s.handleJsonError(r, w, 500, messages.ErrorValidacion)
			return
		}

		var result = make([]Attachment, 0)

		db := database.GetDB()
		rows, err := db.Query("SELECT id, nombre_archivo, titulo FROM adjuntos WHERE trabajo_id = ?", trabajoId)

		if err != nil {
			log.Error("error leyendo adjuntos: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorDatos)
			return
		}

		for rows.Next() {
			var id int
			var filename string
			var title string
			err = rows.Scan(&id, &filename, &title)
			if err != nil {
				log.Error("error escaneando adjuntos: %s", err.Error())
				s.handleJsonError(r, w, 500, messages.ErrorDatos)
				return
			}
			result = append(result, Attachment{Id: id, Title: title, Filename: filename})
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
