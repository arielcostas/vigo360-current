/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package internal

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handlePublicAutorAtom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		autorid := mux.Vars(r)["autorid"]
		var autor, err = s.store.autor.Obtener(autorid)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("no se encontró el autor: %s", err.Error())
				s.handleError(w, 404, messages.ErrorPaginaNoEncontrada)
			} else {
				logger.Error("error recuperando el autor: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}
			return
		}

		pp, err := s.store.publicacion.ListarPorAutor(autorid)
		if err != nil {
			logger.Error("error recuperando publicaciones: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}
		pp = pp.FiltrarPublicas()

		lastUpdate, _ := pp.ObtenerUltimaActualizacion()

		var result bytes.Buffer
		err = t.ExecuteTemplate(&result, "atom.xml", atomParams{
			Dominio:    os.Getenv("DOMAIN"),
			Path:       "/autores/" + autorid + "/atom.xml",
			Titulo:     autor.Nombre,
			Subtitulo:  "Últimas publicaciones escritas por " + autor.Nombre,
			LastUpdate: lastUpdate.Format(time.RFC3339),
			Entries:    pp,
		})
		if err != nil {
			logger.Error("error generando feed: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
		w.Header().Add("Content-Type", "application/atom+xml;charset=UTF-8")
		w.Write(result.Bytes())
	}
}
