package internal

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handleAdminEditWorkAction() http.HandlerFunc {
	type EditPostActionFormInput struct {
		Titulo     string `validate:"required,min=3,max=80"`
		Resumen    string `validate:"required,min=3,max=300"`
		Contenido  string `validate:"required"`
		AltPortada string `validate:"required,min=3,max=300"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		trabajoId := mux.Vars(r)["id"]

		_, err := s.store.trabajo.ObtenerPorId(trabajoId, false)
		if err != nil {
			log.Error("no se encontró el trabajo a editar")
			s.handleError(r, w, 404, messages.ErrorPaginaNoEncontrada)
			return
		}

		if err := r.ParseMultipartForm(26214400); err != nil {
			log.Error("no se pudo extraer datos del formulario: %s", err.Error())
			s.handleError(r, w, 404, messages.ErrorFormulario)
			return
		}

		fi := EditPostActionFormInput{
			Titulo:     r.FormValue("work-titulo"),
			Resumen:    r.FormValue("work-resumen"),
			Contenido:  r.FormValue("work-contenido"),
			AltPortada: r.FormValue("alt_portada"),
		}

		if err := validator.New().Struct(fi); err != nil {
			// TODO: Show actual validation error, and form again
			log.Error("error validando el formulario: %s", err.Error())
			s.handleError(r, w, 404, messages.ErrorValidacion)
			return
		}

		var tx *sql.Tx

		if nt, err := database.GetDB().Begin(); err != nil {
			log.Error("error comenzando transacción: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		} else {
			tx = nt
		}

		query := `UPDATE trabajos SET titulo=?, resumen=?, contenido=?, alt_portada=? WHERE id=?`
		if _, err := tx.Exec(query,
			strings.TrimSpace(fi.Titulo),
			strings.TrimSpace(fi.Resumen),
			strings.TrimSpace(fi.Contenido),
			strings.TrimSpace(fi.AltPortada),
			trabajoId,
		); err != nil {
			e2 := tx.Rollback()
			if e2 != nil {
				s.handleError(r, w, 500, messages.ErrorDatos)
			}
			log.Error("error actualizando trabajo: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		if r.FormValue("publicar") == "on" {
			// TODO: Update this above with the others
			query := `UPDATE trabajos SET fecha_publicacion=NOW() WHERE id=?`
			if _, err := tx.Exec(query, trabajoId); err != nil {
				log.Error("error actualizando fecha de publicación: %s", err.Error())
				s.handleError(r, w, 500, messages.ErrorDatos)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			log.Error("error haciendo commit: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		portada_file, _, err := r.FormFile("portada")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			log.Error("error extrayendo imagen: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorValidacion)
			return
		}

		// Image uploaded
		if !errors.Is(err, http.ErrMissingFile) {
			encodeImagesAndSave(portada_file, trabajoId)
		}

		defer w.WriteHeader(303)
		if r.FormValue("salir") == "true" {
			w.Header().Add("Location", "/admin/works")
		} else {
			w.Header().Add("Location", r.URL.Path)
		}
	}
}
