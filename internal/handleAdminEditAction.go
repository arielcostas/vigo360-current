package internal

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/seo"
)

func (s *Server) handleAdminEditAction() http.HandlerFunc {

	type EditPostActionFormInput struct {
		Titulo     string `validate:"required,min=3,max=80"`
		Resumen    string `validate:"required,min=3,max=300"`
		Contenido  string `validate:"required"`
		AltPortada string `validate:"required,min=3,max=300"`

		SerieId       string
		SeriePosicion string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		publicacionId := mux.Vars(r)["id"]

		_, err := s.store.publicacion.ObtenerPorId(publicacionId, false)
		if err != nil {
			log.Error("no se encontró la publicación a editar")
			s.handleError(w, 404, messages.ErrorPaginaNoEncontrada)
			return
		}

		if err := r.ParseMultipartForm(26214400); err != nil {
			log.Error("no se pudo extraer datos del formulario: %s", err.Error())
			s.handleError(w, 404, messages.ErrorFormulario)
			return
		}

		fi := EditPostActionFormInput{
			Titulo:        r.FormValue("art-titulo"),
			Resumen:       r.FormValue("art-resumen"),
			Contenido:     r.FormValue("art-contenido"),
			AltPortada:    r.FormValue("alt-portada"),
			SerieId:       r.FormValue("serie-id"),
			SeriePosicion: r.FormValue("serie-num"),
		}

		if err := validator.New().Struct(fi); err != nil {
			// TODO: Show actual validation error, and form again
			log.Error("error validando el formulario: %s", err.Error())
			s.handleError(w, 404, messages.ErrorValidacion)
			return
		}

		tags := r.Form["tags"]
		var tx *sql.Tx

		if nt, err := database.GetDB().Begin(); err != nil {
			log.Error("error comenzando transacción: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		} else {
			tx = nt
		}

		if _, err := tx.Exec("DELETE FROM publicaciones_tags WHERE publicacion_id = ?", publicacionId); err != nil {
			tx.Rollback()
			log.Error("error eliminando tags existentes: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		for _, t := range tags {
			if _, err := tx.Exec("INSERT INTO publicaciones_tags (publicacion_id, tag_id) VALUES (?, ?)", publicacionId, t); err != nil {
				tx.Rollback()
				log.Error("error insertando nuevas tags: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
				return
			}
		}

		query := `UPDATE publicaciones SET titulo=?, resumen=?, contenido=?, alt_portada=? WHERE id=?`
		if _, err := tx.Exec(query, fi.Titulo, fi.Resumen, fi.Contenido, fi.AltPortada, publicacionId); err != nil {
			tx.Rollback()
			log.Error("error actualizando publicación: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		if r.FormValue("publicar") == "on" {
			query := `UPDATE publicaciones SET fecha_publicacion=NOW() WHERE id=?`
			if _, err := tx.Exec(query, publicacionId); err != nil {
				log.Error("error actualizando fecha de publicación: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
				return
			}

			var DOMAIN = os.Getenv("DOMAIN")
			var indexnowurls = []string{
				DOMAIN + "/",
				DOMAIN + "/post/" + publicacionId,
			}

			for _, t := range tags {
				indexnowurls = append(indexnowurls, DOMAIN+"/tags/"+t+"/")
			}

			err = seo.BingIndexnowRequest(indexnowurls)
			if err != nil {
				log.Error("error llamando a indexnow: %s", err.Error())
			}
		}

		if fi.SerieId != "" {
			if fi.SeriePosicion == "" {
				fi.SeriePosicion = "1"
			}

			if _, err := tx.Exec(`UPDATE publicaciones SET serie_id = ?, serie_posicion = ? WHERE id = ?`, fi.SerieId, fi.SeriePosicion, publicacionId); err != nil {
				tx.Rollback()
				log.Error("error guardando serie: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
				return
			}
		}

		if err := tx.Commit(); err != nil {
			log.Error("error haciendo commit: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		portada_file, _, err := r.FormFile("portada")
		if err != nil && !errors.Is(err, http.ErrMissingFile) {
			log.Error("error extrayendo imagen: %s", err.Error())
			s.handleError(w, 500, messages.ErrorValidacion)
			return
		}

		// Image uploaded
		if !errors.Is(err, http.ErrMissingFile) {
			encodeImagesAndSave(portada_file, publicacionId)
		}

		defer w.WriteHeader(303)
		if r.FormValue("salir") == "true" {
			w.Header().Add("Location", "/admin/post")
		} else {
			w.Header().Add("Location", r.URL.Path)
		}
	}
}

func encodeImagesAndSave(portada_file io.Reader, publicacion_id string) {
	uppath := os.Getenv("UPLOAD_PATH")
	var err error
	logger := logger.NewLogger("encodeImagesAndSave " + publicacion_id)

	var portadaJpg, portadaWebp bytes.Buffer
	if pj, pw, e2 := generateImagesFromImage(portada_file); errors.Is(e2, ErrImageFormatError) {
		logger.Error("error procesando imágenes: %s", err.Error())
		return
	} else if err != nil {
		logger.Error("error procesando imágenes: %s", err.Error())
		return
	} else {
		portadaJpg = pj
		portadaWebp = pw
	}

	if e2 := os.WriteFile(uppath+"/thumb/"+publicacion_id+".jpg", portadaJpg.Bytes(), os.ModePerm); e2 != nil {
		logger.Error("error guardando imagen jpg: %s", err.Error())
		return
	}

	if e2 := os.WriteFile(uppath+"/images/"+publicacion_id+".webp", portadaWebp.Bytes(), os.ModePerm); e2 != nil {
		logger.Error("error guardando imagen webp: %s", err.Error())
		return
	}
}
