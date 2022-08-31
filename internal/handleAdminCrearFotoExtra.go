// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/thanhpk/randstr"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handleAdminCrearFotoExtra() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		uploadPath := os.Getenv("UPLOAD_PATH")

		articuloId := r.FormValue("articulo")
		if articuloId == "" {
			logger.Error("el id del artículo no puede estar vacío: %s")
			s.handleJsonError(w, 500, messages.ErrorFormulario)
			return
		}

		file, _, err := r.FormFile("foto")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			logger.Error("no se ha subido ninguna imagen: %s", err.Error())
			s.handleJsonError(w, 500, messages.ErrorFormulario)
			return
		}

		photoBytes, err := io.ReadAll(file)
		if err != nil {
			logger.Error("no se pudo extraer la imagen del formulario: %s", err.Error())
			s.handleJsonError(w, 500, messages.ErrorFormulario)
			return
		}

		image, err := imagenDesdeMime(photoBytes)
		if err != nil {
			logger.Error("error extrayendo el tipo MIME de la imagen: %s", err.Error())
			s.handleJsonError(w, 500, messages.ErrorFormulario)
			return
		}

		var fotoEscribir bytes.Buffer
		err = jpeg.Encode(&fotoEscribir, image, &jpeg.Options{Quality: 95})
		if err != nil {
			logger.Error("error codificando la imagen: %s", err.Error())
			s.handleJsonError(w, 500, messages.ErrorDatos)
			return
		}

		var salt = randstr.String(5)
		var imagePath = fmt.Sprintf("%s/extra/%s-%s.jpg", uploadPath, articuloId, salt)
		err = os.WriteFile(imagePath, fotoEscribir.Bytes(), 0o644)
		if err != nil {
			logger.Error("error escribiendo imagen a %s: %s", imagePath, err.Error())
			s.handleJsonError(w, 500, messages.ErrorRender)
		}
	}
}
