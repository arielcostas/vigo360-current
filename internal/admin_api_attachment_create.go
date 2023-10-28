package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/thanhpk/randstr"
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

		// filename = fileheader.Filename all as lowercase, separated by underscores and up to 40 characters long
		var filename = strings.ToLower(fileheader.Filename)
		filename = strings.ReplaceAll(filename, " ", "_")
		filename = strings.ReplaceAll(filename, "á", "a")
		filename = strings.ReplaceAll(filename, "é", "e")
		filename = strings.ReplaceAll(filename, "í", "i")
		filename = strings.ReplaceAll(filename, "ó", "o")
		filename = strings.ReplaceAll(filename, "ú", "u")
		filename = strings.ReplaceAll(filename, "ñ", "n")
		filename = filename[:40]
		var salt = randstr.String(3) + "_"
		filename = salt + filename

		_, err = tx.Exec(
			"INSERT INTO adjuntos (trabajo_id, nombre_archivo, titulo) VALUES (?, ?, ?)",
			trabajoId, filename, titulo,
		)
		if err != nil {
			log.Error("error guardando a bbddd: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		var imagePath = fmt.Sprintf("%s/papers/%s", uploadPath, filename)
		err = os.WriteFile(imagePath, fileBytes, 0o644)
		if err != nil {
			_ = tx.Rollback()
			log.Error("error escribiendo imagen a %s: %s", imagePath, err.Error())
			s.handleJsonError(r, w, 500, messages.ErrorRender)
			return
		} else {
			_ = tx.Commit()
		}

		w.WriteHeader(201)
		w.Write([]byte("{}\n"))
		w.Header().Add("Content-Type", "application/json")
	}
}
