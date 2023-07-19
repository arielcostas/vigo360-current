package internal

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/nfnt/resize"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
)

func (s *Server) handleAdminPerfilEdit() http.HandlerFunc {
	type entradaFormulario struct {
		Nombre     string `validate:"required,min=3,max=80"`
		Biografia  string `validate:"required,min=3,max=2000"`
		Web_titulo string `validate:"max=80"`
		Web_url    string `validate:"max=80"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)

		_, err := s.store.autor.Obtener(sess.Autor_id)
		if err != nil {
			logger.Error("no se encontró el autor a editar")
			s.handleError(r, w, 404, messages.ErrorPaginaNoEncontrada)
			return
		}

		if err := r.ParseMultipartForm(26214400); err != nil {
			logger.Error("no se pudo extraer datos del formulario: %s", err.Error())
			s.handleError(r, w, 404, messages.ErrorFormulario)
			return
		}

		fi := entradaFormulario{
			Nombre:     r.FormValue("nombre"),
			Biografia:  r.FormValue("biografia"),
			Web_titulo: r.FormValue("web-titulo"),
			Web_url:    r.FormValue("web-url"),
		}

		if err := validator.New().Struct(fi); err != nil {
			// TODO: Show actual validation error, and form again
			logger.Error("error validando el formulario: %s", err.Error())
			s.handleError(r, w, 404, messages.ErrorValidacion)
			return
		}

		query := `UPDATE autores SET nombre=?, biografia=?, web_titulo=?, web_url=? WHERE id=?`
		if _, err := database.GetDB().Exec(query, fi.Nombre, fi.Biografia, fi.Web_titulo, fi.Web_url, sess.Autor_id); err != nil {
			logger.Error("error actualizando autor: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		perfil_file, _, err := r.FormFile("perfil")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			logger.Error("error extrayendo nueva foto de perfil: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorValidacion)
			return
		}

		if !errors.Is(err, http.ErrMissingFile) {
			uppath := os.Getenv("UPLOAD_PATH")

			photoBytes, err := io.ReadAll(perfil_file)
			if err != nil {
				logger.Error("no se pudo extraer la imagen del formulario: %s", err.Error())
				s.handleJsonError(r, w, 500, messages.ErrorFormulario)
				return
			}

			image, err := imagenDesdeMime(photoBytes)
			if err != nil {
				logger.Error("error extrayendo el tipo MIME de la imagen: %s", err.Error())
				s.handleJsonError(r, w, 500, messages.ErrorFormulario)
				return
			}

			image = resize.Resize(256, 256, image, resize.Bicubic)

			var fotoEscribir bytes.Buffer
			err = jpeg.Encode(&fotoEscribir, image, &jpeg.Options{Quality: 95})
			if err != nil {
				logger.Error("error codificando la imagen: %s", err.Error())
				s.handleJsonError(r, w, 500, messages.ErrorDatos)
				return
			}

			var imagePath = fmt.Sprintf("%s/profile/%s.jpg", uppath, sess.Autor_id)
			err = os.WriteFile(imagePath, fotoEscribir.Bytes(), 0o644)
			if err != nil {
				logger.Error("error escribiendo imagen a %s: %s", imagePath, err.Error())
				s.handleJsonError(r, w, 500, messages.ErrorRender)
			}
		}

		defer w.WriteHeader(303)
		w.Header().Add("Location", r.URL.Path)
	}
}
