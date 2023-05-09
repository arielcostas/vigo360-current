package internal

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handleAdminCreateSeries() http.HandlerFunc {
	type entrada struct {
		Titulo string `validate:"required,min=1,max=40"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		if err := r.ParseForm(); err != nil {
			logger.Error("error obteniendo datos del formulario: %s", err.Error())
			s.handleError(w, 500, messages.ErrorFormulario)
			return
		}

		fi := entrada{}
		fi.Titulo = r.FormValue("titulo")

		if err := validator.New().Struct(fi); err != nil {
			logger.Error("error validando parámetros: %s", err.Error())
			s.handleError(w, 500, messages.ErrorValidacion)
			return
		}

		id := strings.ToLower(strings.TrimSpace(fi.Titulo))
		id = strings.ReplaceAll(id, " ", "-")

		if _, err := database.GetDB().Exec(`INSERT INTO series VALUES (?, ?)`, id, fi.Titulo); err != nil {
			logger.Error("error creando serie: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		w.Header().Add("Location", "/admin/series")
		w.WriteHeader(303)
	}
}
