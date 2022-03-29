/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"errors"
	"net/http"
	"os"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
)

func EditPostPage(w http.ResponseWriter, r *http.Request) {
	// TODO: Check author is same as session
	verifyLogin(w, r)
	post_id := mux.Vars(r)["id"]

	post := PostEditar{}

	err := db.QueryRowx(`SELECT id, titulo, resumen, contenido, alt_portada, (fecha_publicacion is not null && fecha_publicacion < NOW()) as publicado, serie_id, serie_posicion FROM publicaciones WHERE id = ?;`, post_id).StructScan(&post)

	// TODO: Proper error handling
	if err != nil {
		logger.Error("[editor]: error getting article from database: %s", err.Error())
		w.WriteHeader(500)
		_, err2 := w.Write([]byte("error buscando el artículo en la base de datos"))
		if err2 != nil {
			logger.Error("[post] error showing error message: %s", err2.Error())
		}
		return
	}

	series := []Serie{}
	err = db.Select(&series, `SELECT * FROM series;`)

	// TODO: Proper error handling
	if err != nil {
		logger.Error("[editor] error getting article from database: %s", err.Error())
		w.WriteHeader(500)
		_, err2 := w.Write([]byte("error buscando el artículo en la base de datos"))
		if err2 != nil {
			logger.Error("[post] error showing error message: %s", err2.Error())
		}
		return
	}

	tags := []Tag{}
	err = db.Select(&tags, `SELECT id, nombre, (SELECT tag_id FROM publicaciones_tags pt WHERE pt.publicacion_id = ? AND pt.tag_id = id) IS NOT NULL as seleccionada FROM tags`, post_id)
	if err != nil {
		logger.Error("[editor] error fetching tags: %s", err.Error())
		w.WriteHeader(500)
		w.Write([]byte("error obteniendo datos"))
		return
	}

	err = t.ExecuteTemplate(w, "post-id.html", struct {
		Post   PostEditar
		Series []Serie
		Tags   []Tag
	}{
		Post:   post,
		Series: series,
		Tags:   tags,
	})
	if err != nil {
		logger.Error("[editor-postid] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
}

func EditPostAction(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	publicacion_id := mux.Vars(r)["id"]
	err := r.ParseMultipartForm(26214400) // 25 MB

	if err != nil {
		w.WriteHeader(500)
		logger.Error("error parsing multipart form: %s", err.Error())
		InternalServerErrorHandler(w, r)
	}

	art_titulo := r.FormValue("art-titulo")
	art_resumen := r.FormValue("art-resumen")
	art_contenido := r.FormValue("art-contenido")
	alt_portada := r.FormValue("alt-portada")
	art_publicar := r.FormValue("publicar")

	serie_id := r.FormValue("serie-id")
	serie_posicion := r.FormValue("serie-num")

	tags := r.Form["tags"]

	tx, err := db.Begin()
	if err != nil {
		logger.Error("[editor] error beginning transaction: %s", err.Error())
		w.WriteHeader(500)
		w.Write([]byte("error guardando datos"))
		return
	}
	tx.Exec("DELETE FROM tags WHERE publicacion_id = ?", publicacion_id)
	for _, t := range tags {
		_, err = tx.Exec("INSERT INTO publicaciones_tags (publicacion_id, tag_id) VALUES (?, ?)", publicacion_id, t)
		if err != nil {
			logger.Error("[editor] error adding tag to post %s: %s", publicacion_id, err.Error())
			w.WriteHeader(400)
			w.Write([]byte("error guardando datos"))
			return
		}
	}

	// TODO: Make all edits in one transaction, add locks
	err = tx.Commit()
	if err != nil {
		logger.Error("[editor] error commiting tag editing %s: %s", publicacion_id, err.Error())
		w.WriteHeader(400)
		w.Write([]byte("error guardando datos"))
		return
	}

	// TODO: Proper error page
	if !validarTitulo(art_titulo) {
		w.WriteHeader(400)
		_, err2 := w.Write([]byte("El título debe contener entre 3 y 80 caracteres"))
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
		return
	}

	if !validarResumen(art_resumen) {
		w.WriteHeader(400)
		_, err2 := w.Write([]byte("El resumen debe contener entre 3 y 300 caracteres"))
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
		return
	}

	if !validarContenido(art_contenido) {
		w.WriteHeader(400)
		_, err2 := w.Write([]byte("El contenido del artículo no puede estar vacío"))
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
		return
	}

	// TODO: Refactor this utter piece of crap
	query := `UPDATE publicaciones SET titulo=?, resumen=?, contenido=?, alt_portada=?, serie_id=?, serie_posicion=?`
	if art_publicar == "on" {
		query += `, fecha_publicacion = NOW()`
	}

	// If serie is unselected but posicion is not deleted, it will be saved, even though it doesn't make sense
	if serie_id == "" {
		serie_posicion = ""
	}
	_, err = db.Exec(query+` WHERE id=?`, art_titulo, art_resumen, art_contenido, alt_portada, NewNullString(serie_id), NewNullString(serie_posicion), publicacion_id)

	// TODO: Proper error page
	if err != nil {
		logger.Error("[editor] error saving edited post to database: %s", err.Error())
		w.WriteHeader(400)
		_, err2 := w.Write([]byte("error guardando cambios a la base de datos"))
		if err2 != nil {
			logger.Error("[editor] error reverting database: %s", err2.Error())
		}
	}

	logger.Information("[editor] updated post %s", publicacion_id)

	// image processing
	portada_file, _, err := r.FormFile("portada")

	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		logger.Error("[editor] unexpected error extracting image: %s:", err.Error())
		return
	}

	if !errors.Is(err, http.ErrMissingFile) {
		portadaJpg, portadaWebp, err := generateImagesFromImage(portada_file)
		if errors.Is(err, InvalidImageFormatError) {
			logger.Error("[editor] user uploaded image with invalid mime")
			w.Write([]byte("La imagen subida no tiene un formato válido"))
			return
		} else if err != nil {
			logger.Error("unexpected error generating images: %s", err)
			return
		}

		err = os.WriteFile(os.Getenv("UPLOAD_PATH")+"/thumb/"+publicacion_id+".jpg", portadaJpg.Bytes(), os.ModePerm)
		if err != nil {
			logger.Error("[editor] error writing jpg image: %s", err)
		}

		err = os.WriteFile(os.Getenv("UPLOAD_PATH")+"/images/"+publicacion_id+".webp", portadaWebp.Bytes(), os.ModePerm)
		if err != nil {
			logger.Error("[editor] error writing webp file: %s", err)
		}
	}

	w.Header().Add("Location", "/admin/post")
	defer w.WriteHeader(303)
}
