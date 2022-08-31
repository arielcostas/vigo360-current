// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
)

func (s *Server) handle_api_listar_comentarios(w http.ResponseWriter, r *http.Request) {
	var comentarios, err = s.store.comentario.ListarPorEstado(models.ESTADO_PENDIENTE)
	logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))

	if err != nil {
		s.handleJsonError(w, 500, messages.ErrorDatos)
		logger.Error("cannot get comentarios: " + err.Error())
		return
	}

	if len(comentarios) == 0 {
		fmt.Fprintf(w, "[]")
		return
	}

	var salida []byte
	salida, err = json.Marshal(comentarios)

	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "{ error: \"Error de respuesta\" }")
		logger.Error("Error de json: " + err.Error())
		return
	}

	w.Write(salida)
}
