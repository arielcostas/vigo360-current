/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handleAdminListarFotoExtra() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		uploadPath := os.Getenv("UPLOAD_PATH")

		articuloId := r.URL.Query().Get("articulo")
		if articuloId == "" {
			logger.Error("no se especificó un articuloId")
			s.handleJsonError(w, 500, messages.ErrorValidacion)
			return
		}

		files, err := os.ReadDir(uploadPath + "/extra")
		if err != nil {
			logger.Error("error leyendo carpeta %s: %s", uploadPath+"/extra", err.Error())
			s.handleJsonError(w, 500, messages.ErrorDatos)
			return
		}

		var result = make([]string, 0)

		for _, de := range files {
			if de.IsDir() {
				continue
			}

			name := de.Name()
			if strings.HasPrefix(name, articuloId) && strings.HasSuffix(name, ".jpg") {
				result = append(result, name)
			}
		}

		if len(result) < 1 {
			logger.Error("no se encontró ningún resultado")
			s.handleJsonError(w, 404, "Ningún resultado encontrado")
			return
		}

		resbytes, err := json.MarshalIndent(result, "", "\t")
		if err != nil {
			logger.Error("error escribiendo json de respuesta: %s", err.Error())
			s.handleJsonError(w, 500, messages.ErrorRender)
			return
		}
		w.Write(resbytes)
	}
}
