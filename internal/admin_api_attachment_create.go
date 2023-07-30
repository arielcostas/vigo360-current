package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) adminApiAttachmentCreate() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		uploadPath := os.Getenv("UPLOAD_PATH")

		trabajoId := r.FormValue("trabajo")
		if trabajoId == "" {
			log.Error("el id del trabajo no puede estar vacío: %s")
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		titulo := r.FormValue("titulo")
		if trabajoId == "" {
			log.Error("el titulo del adjunto no puede estar vacío: %s")
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		file, fileheader, err := r.FormFile("file")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			log.Error("no se ha subido ningun archivo: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			log.Error("no se pudo extraer el archivo del formulario: %s", err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorFormulario)
			return
		}

		var tx *sql.Tx
		if nt, err := database.GetDB().Begin(); err != nil {
			log.Error("error iniciando transacción: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		} else {
			tx = nt
		}

		_, err = tx.Exec(
			"INSERT INTO adjuntos (trabajo_id, nombre_archivo, titulo) VALUES (?, ?, ?)",
			trabajoId, fileheader.Filename, titulo,
		)
		if err != nil {
			return
		}

		var imagePath = fmt.Sprintf("%s/trabajos/%s", uploadPath, fileheader.Filename)
		err = os.WriteFile(imagePath, fileBytes, 0o644)
		if err != nil {
			_ = tx.Rollback()
			log.Error("error escribiendo imagen a %s: %s", imagePath, err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorRender)
		} else {
			_ = tx.Commit()
		}

		w.WriteHeader(201)
	}
}
