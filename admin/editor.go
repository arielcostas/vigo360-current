package admin

import (
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/gorilla/mux"
)

func EditPostPage(w http.ResponseWriter, r *http.Request) {
	// TODO Check author is same as session
	verifyLogin(w, r)
	post_id := mux.Vars(r)["id"]
	post := PostEditar{}

	err := db.QueryRowx(`SELECT id, titulo, resumen, contenido, alt_portada, (fecha_publicacion is not null && fecha_publicacion < NOW()) as publicado FROM publicaciones WHERE id = ?;`, post_id).StructScan(&post)

	// TODO Proper error handling
	if err != nil {
		logger.Error("[editor]: error getting article from database: %s", err.Error())
		w.WriteHeader(500)
		_, err2 := w.Write([]byte("error buscando el artículo en la base de datos"))
		if err2 != nil {
			logger.Error("[post] error showing error message: %s", err2.Error())
		}
		return
	}

	err = t.ExecuteTemplate(w, "post-id.html", struct {
		Post PostEditar
	}{Post: post})
	if err != nil {
		logger.Error("[editor-postid] error rendering template: %s", err.Error())
		InternalServerErrorHandler(w, r)
		return
	}
}

func EditPostAction(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	post_id := mux.Vars(r)["id"]
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

	// TODO Proper error page
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

	query := `UPDATE publicaciones SET titulo=?, resumen=?, contenido=?, alt_portada=?`
	if art_publicar == "on" {
		query += `, fecha_publicacion = NOW()`
	}
	_, err = db.Exec(query+` WHERE id=?`, art_titulo, art_resumen, art_contenido, alt_portada, post_id)

	// TODO Proper error page
	if err != nil {
		logger.Error("error saving edited post to database: %s", err.Error())
		w.WriteHeader(400)
		_, err2 := w.Write([]byte("error guardando cambios a la base de datos"))
		if err2 != nil {
			logger.Error("[post] error reverting database: %s", err2.Error())
		}
	}

	logger.Information("[editor] updated post %s", post_id)

	// image processing
	portada_file, portada_mime, err := r.FormFile("portada")

	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		log.Fatalln("error imagen 71:" + err.Error())
	}

	if !errors.Is(err, http.ErrMissingFile) {

		portadaJpg, portadaWebp := generateImagesFromImage(portada_file, portada_mime)

		file, err := os.Create(os.Getenv("UPLOAD_PATH") + "/thumb/" + post_id + ".jpg")
		if err != nil {
			log.Fatalf("error opening file for writing: %s", err)
		}
		_, err = file.Write(portadaJpg.Bytes())
		if err != nil {
			log.Fatalf("error writing file: %s", err)
		}

		file, err = os.Create(os.Getenv("UPLOAD_PATH") + "/images/" + post_id + ".webp")
		if err != nil {
			log.Fatalf("error opening file for writing: %s", err)
		}
		_, err = file.Write(portadaWebp.Bytes())
		if err != nil {
			log.Fatalf("error writing file: %s", err)
		}
	}

	// If not, redirection shows the old value - Chapuza xD
	time.Sleep(100 * time.Millisecond)
	w.Header().Add("Location", "/admin/post")
	w.WriteHeader(303)
}
