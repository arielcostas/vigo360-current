/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"vigo360.es/new/internal/logger"
)

func postEditor(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}
	sess, err := getSession(sc.Value)
	if err != nil {
		logger.Notice("unauthenticated user tried to access this page")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}
	post_id := mux.Vars(r)["id"]

	// TODO: Check author is same as session
	post := PostEditar{}

	err = db.QueryRowx(`SELECT id, titulo, resumen, contenido, alt_portada, (fecha_publicacion is not null && fecha_publicacion < NOW()) as publicado, serie_id, serie_posicion FROM publicaciones WHERE id = ?;`, post_id).StructScan(&post)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &appError{Error: err, Message: "trying to edit not-found article",
				Response: "No se encontró el artículo a editar.", Status: 404}
		}

		return newDatabaseReadAppError(err, "post")
	}

	series := []Serie{}
	err = db.Select(&series, `SELECT * FROM series;`)

	if err != nil {
		return newDatabaseReadAppError(err, "series")
	}

	tags := []Tag{}
	err = db.Select(&tags, `SELECT id, nombre, (SELECT tag_id FROM publicaciones_tags pt WHERE pt.publicacion_id = ? AND pt.tag_id = id) IS NOT NULL as seleccionada FROM tags ORDER BY nombre ASC`, post_id)
	if err != nil {
		return newDatabaseReadAppError(err, "tags")
	}

	err = t.ExecuteTemplate(w, "post-id.html", struct {
		Post    PostEditar
		Series  []Serie
		Tags    []Tag
		Session Session
	}{
		Post:    post,
		Series:  series,
		Tags:    tags,
		Session: sess,
	})
	if err != nil {
		return newTemplateRenderingAppError(err)
	}
	return nil
}

type EditPostActionFormInput struct {
	Titulo      string `validate:"required,min=3,max=80"`
	Resumen     string `validate:"required,min=3,max=300"`
	Contenido   string `validate:"required"`
	Alt_portada string `validate:"required,min=3,max=300"`

	Serie_id       string
	Serie_posicion string
}

func editPost(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}
	_, err = getSession(sc.Value)
	if err != nil {
		logger.Notice("unauthenticated user tried to access this page")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return nil
	}

	publicacion_id := mux.Vars(r)["id"]

	// TODO: Check post exists before parsing form
	if err := r.ParseMultipartForm(26214400); err != nil {
		return &appError{Error: err, Message: "error parsing multipart form",
			Response: "Hubo un error recibiendo los datos del formulario", Status: 500}
	}

	fi := EditPostActionFormInput{
		Titulo:         r.FormValue("art-titulo"),
		Resumen:        r.FormValue("art-resumen"),
		Contenido:      r.FormValue("art-contenido"),
		Alt_portada:    r.FormValue("alt-portada"),
		Serie_id:       r.FormValue("serie-id"),
		Serie_posicion: r.FormValue("serie-num"),
	}

	if err := validator.New().Struct(fi); err != nil {
		// TODO: Show actual validation error, and form again
		return &appError{Error: err, Message: "form data is not valid",
			Response: "Error de validación del formulario", Status: 400}
	}

	tags := r.Form["tags"]
	var tx *sql.Tx

	if nt, err := db.Begin(); err != nil {
		return &appError{Error: err, Message: "error beginning transaction",
			Response: "Hubo un error guardando los cambios", Status: 500}
	} else {
		tx = nt
	}

	if _, err := tx.Exec("DELETE FROM publicaciones_tags WHERE publicacion_id = ?", publicacion_id); err != nil {
		return &appError{Error: err, Message: "error deleting tags from post " + publicacion_id,
			Response: "Hubo un error guardando los cambios", Status: 500}
	}

	for _, t := range tags {
		if _, err := tx.Exec("INSERT INTO publicaciones_tags (publicacion_id, tag_id) VALUES (?, ?)", publicacion_id, t); err != nil {
			return &appError{Error: err, Message: "error saving tag for post " + publicacion_id,
				Response: "Hubo un error guardando los cambios", Status: 500}
		}
	}

	query := `UPDATE publicaciones SET titulo=?, resumen=?, contenido=?, alt_portada=? WHERE id=?`
	if _, err := tx.Exec(query, fi.Titulo, fi.Resumen, fi.Contenido, fi.Alt_portada, publicacion_id); err != nil {
		tx.Rollback()
		return &appError{Error: err, Message: "error saving new post data for " + publicacion_id,
			Response: "Hubo un error guardando los cambios", Status: 500}
	}

	if r.FormValue("publicar") == "on" {
		query := `UPDATE publicaciones SET fecha_publicacion=NOW() WHERE id=?`
		if _, err := tx.Exec(query, publicacion_id); err != nil {
			return &appError{Error: err, Message: "error saving new post data for " + publicacion_id,
				Response: "Hubo un error guardando los cambios", Status: 500}
		}
	}

	if fi.Serie_id != "" {
		if fi.Serie_posicion == "" {
			fi.Serie_posicion = "1"
		}

		if _, err := tx.Exec(`UPDATE publicaciones SET serie_id = ?, serie_posicion = ? WHERE id = ?`, fi.Serie_id, fi.Serie_posicion, publicacion_id); err != nil {
			tx.Rollback()
			return &appError{Error: err, Message: "error adding series for " + publicacion_id,
				Response: "Hubo un error guardando los cambios", Status: 500}
		}
	}

	if err := tx.Commit(); err != nil {
		return &appError{Error: err, Message: "error committing transaction for " + publicacion_id,
			Response: "Hubo un error guardando los cambios", Status: 500}
	}

	portada_file, _, err := r.FormFile("portada")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return &appError{Error: err, Message: "error extracting image from form " + publicacion_id,
			Response: "Hubo un error modificando la imagen. El resto de datos se han guardado.", Status: 400}
	}

	// Image uploaded
	if !errors.Is(err, http.ErrMissingFile) {
		uppath := os.Getenv("UPLOAD_PATH")
		if uppath == "" {
			return &appError{Message: "UPLOAD_PATH is not set, images cannot be stored",
				Response: "El servidor no puede guardar imágenes en este momento."}
		}

		var portadaJpg, portadaWebp bytes.Buffer
		if pj, pw, e2 := generateImagesFromImage(portada_file); errors.Is(e2, ErrImageFormatError) {
			return &appError{Error: e2, Message: "error processing uploaded image",
				Response: "La imagen subida no tiene un formato válido. El resto de datos se han guardado.", Status: 400}
		} else if err != nil {
			return &appError{Error: e2, Message: "unexpected error processing uploaded image",
				Response: "Error inesperado procesando la imagen. El resto de datos se han guardado.", Status: 500}
		} else {
			portadaJpg = pj
			portadaWebp = pw
		}

		if e2 := os.WriteFile(uppath+"/thumb/"+publicacion_id+".jpg", portadaJpg.Bytes(), os.ModePerm); e2 != nil {
			return &appError{Error: e2, Message: "error saving jpg image for " + publicacion_id,
				Response: "Error guardando la imagen. El resto de cambios se han guradado.", Status: 500}
		}

		if e2 := os.WriteFile(uppath+"/images/"+publicacion_id+".webp", portadaWebp.Bytes(), os.ModePerm); e2 != nil {
			return &appError{Error: e2, Message: "error saving webp image for " + publicacion_id,
				Response: "Error guardando la imagen. El resto de cambios se han guradado.", Status: 500}
		}
	}

	w.Header().Add("Location", "/admin/post")
	defer w.WriteHeader(303)
	return nil
}
