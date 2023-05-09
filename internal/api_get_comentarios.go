package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
)

func (s *Server) handle_api_listar_comentarios(w http.ResponseWriter, r *http.Request) {
	var req_raw_estado = r.URL.Query().Get("estado")
	var req_raw_estado_int, err = strconv.Atoi(req_raw_estado)

	var estado models.EstadoComentario
	if err != nil && req_raw_estado != "" {
		s.handleJsonError(w, 400, "estado debe ser un número entre 1 y 3")
		return
	} else if req_raw_estado_int < 1 || req_raw_estado_int > 3 {
		estado = models.ESTADO_PENDIENTE
	} else {
		estado = models.EstadoComentario(req_raw_estado_int)
	}

	var comentarios []models.Comentario
	comentarios, err = s.store.comentario.ListarPorEstado(estado)
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
		s.handleJsonError(w, 500, messages.ErrorRender)
		logger.Error("Error de json: " + err.Error())
		return
	}

	w.Write(salida)
}
