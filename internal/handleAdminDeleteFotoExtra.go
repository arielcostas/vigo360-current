package internal

import (
	"errors"
	"net/http"
	"os"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handleAdminDeleteFotoExtra() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))

		uploadPath := os.Getenv("UPLOAD_PATH")

		fotoId := r.URL.Query().Get("foto")
		if fotoId == "" {
			logger.Error("falta el parámetro foto en la URL")
			s.handleJsonError(w, 500, messages.ErrorFormulario)
			return
		}

		/*
			Comprobar que la fotografía existe, y si es así eliminarla
		*/
		if f, err := os.Stat(uploadPath + "/extra/" + fotoId); err == nil {
			e2 := os.Remove(uploadPath + "/extra/" + f.Name())
			if e2 != nil {
				logger.Error("error borrando fotografía: %s", err.Error())
				s.handleJsonError(w, 500, "Error borrando fotografía")
				return
			}
			w.Write([]byte("{ \"error\": false }"))
		} else if errors.Is(err, os.ErrNotExist) {
			logger.Error("la fotografía no existe: %s", err.Error())
			s.handleJsonError(w, 404, messages.ErrorNoResultados)
		} else {
			logger.Error("error encontrando archivo a borrar: %s", err.Error())
			s.handleJsonError(w, 500, messages.ErrorNoResultados)
		}
	}
}
