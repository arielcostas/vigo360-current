// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package internal

import (
	"database/sql"
	_ "embed"
	"net/http"
	"os"
	"regexp"

	"github.com/go-playground/validator/v10"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
)

//TODO: Replace this

//go:embed extra/default.jpg
var defaultImageJPG []byte

//go:embed extra/default.webp
var defaultImageWebp []byte

func (s *Server) handleAdminCreatePost() http.HandlerFunc {
	type entrada struct {
		Titulo string `validate:"required,min=3,max=80"`
	}

	var postIdRegexp = regexp.MustCompile(`^[A-Za-z0-9\-\_]{3,40}$`)
	var ValidatePostId = func(id string) bool {
		return postIdRegexp.MatchString(id)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)

		var art_autor = sess.Autor_id

		if err := r.ParseForm(); err != nil {
			logger.Error("error leyendo datos de formulario: %s", err.Error())
			s.handleError(w, 500, messages.ErrorFormulario)
			return
		}

		fi := entrada{
			Titulo: r.FormValue("art-titulo"),
		}

		var id = r.FormValue("art-id")

		if !ValidatePostId(id) {
			logger.Error("error validando id")
			s.handleError(w, 500, messages.ErrorValidacion)
			return
		}

		if err := validator.New().Struct(fi); err != nil {
			// TODO: Mostrar de nuevo formulario con los errores
			logger.Error("error validando título: %s", err.Error())
			s.handleError(w, 500, messages.ErrorValidacion)
			return
		}

		var tx *sql.Tx
		if nt, err := database.GetDB().Begin(); err != nil {
			logger.Error("error iniciando transacción: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		} else {
			tx = nt
		}

		q := "INSERT INTO publicaciones(id, titulo, alt_portada, resumen, contenido, autor_id) VALUES (?, ?, 'CAMBIAME','', '', ?);"
		if _, err := tx.Exec(q, id, fi.Titulo, art_autor); err != nil {
			tx.Rollback()
			logger.Error("error creando artículo: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		photopath := os.Getenv("UPLOAD_PATH")
		if err := os.WriteFile(photopath+"/images/"+id+".webp", defaultImageWebp, 0o644); err != nil {
			tx.Rollback()
			logger.Error("error escribiendo imagen webp: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		if err := os.WriteFile(photopath+"/thumb/"+id+".jpg", defaultImageJPG, 0o644); err != nil {
			tx.Rollback()
			logger.Error("error escribiendo imagen jpg: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			logger.Error("error ejecutando transacción: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		w.Header().Add("Location", "/admin/post/"+id)
		w.WriteHeader(303)
	}
}
