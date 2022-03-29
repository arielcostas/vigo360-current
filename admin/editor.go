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
	"github.com/go-playground/validator/v10"
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

type EditPostActionFormInput struct {
	Titulo      string `validate:"required,min=3,max=80"`
	Resumen     string `validate:"required,min=3,max=300"`
	Contenido   string `validate:"required"`
	Alt_portada string `validate:"required,min=3,max=300"`

	Serie_id       string
	Serie_posicion string
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

	fi := EditPostActionFormInput{}

	fi.Titulo = r.FormValue("art-titulo")
	fi.Resumen = r.FormValue("art-resumen")
	fi.Contenido = r.FormValue("art-contenido")
	fi.Alt_portada = r.FormValue("alt-portada")
	fi.Serie_id = r.FormValue("serie-id")
	fi.Serie_posicion = r.FormValue("serie-num")

	err = validator.New().Struct(fi)
	if err != nil {
		logger.Error("[editor] error validating form input: %s", err.Error())
		w.WriteHeader(400)
		w.Write([]byte("Error de validación"))
		return
	}

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

	// TODO: Refactor this utter piece of crap
	query := `UPDATE publicaciones SET titulo=?, resumen=?, contenido=?, alt_portada=?, serie_id=?, serie_posicion=?`
	if r.FormValue("publicar") == "on" {
		query += `, fecha_publicacion = NOW()`
	}

	// If serie is unselected but posicion is not deleted, it will be saved, even though it doesn't make sense
	if fi.Serie_id == "" {
		fi.Serie_posicion = ""
	}
	_, err = db.Exec(query+` WHERE id=?`, fi.Titulo, fi.Resumen, fi.Contenido, fi.Alt_portada, NewNullString(fi.Serie_id), NewNullString(fi.Serie_posicion), publicacion_id)

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
