package internal

import (
	"errors"
	"net/http"
	"os"
	"vigo360.es/new/internal/database"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) adminApiAttachmentDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))

		uploadPath := os.Getenv("UPLOAD_PATH")

		fotoId := r.URL.Query().Get("id")
		if fotoId == "" {
			log.Error("falta el parámetro foto en la URL")
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		db := database.GetDB()

		row := db.QueryRow("SELECT nombre_archivo FROM adjuntos WHERE id = ?", fotoId)
		var nombreArchivo string
		if err := row.Scan(&nombreArchivo); err != nil {
			log.Error("error escaneando nombre de archivo: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorNoResultados)
			return
		}

		_, _ = db.Exec("DELETE FROM adjuntos WHERE id = ? LIMIT 1", fotoId)

		if f, err := os.Stat(uploadPath + "/papers/" + nombreArchivo); err == nil {
			e2 := os.Remove(uploadPath + "/papers/" + f.Name())
			if e2 != nil {
				log.Error("error borrando fotografía: %s", err.Error())
				s.handleJsonError(r, w, 500, "Error borrando fotografía")
				return
			}
			w.WriteHeader(204)
		} else if errors.Is(err, os.ErrNotExist) {
			log.Error("attachment doesn't exist: %s", err.Error())
			s.handleJsonError(r, w, 404, messages.ErrorNoResultados)
		} else {
			log.Error("error finding file to remove: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorNoResultados)
		}
	}
}
