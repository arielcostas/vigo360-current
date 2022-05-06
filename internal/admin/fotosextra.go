/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package admin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/thanhpk/randstr"
	"vigo360.es/new/internal/logger"
)

func listarFotosExtra(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		return LoginRequiredAppError
	}
	_, err = getSession(sc.Value)
	if err != nil {
		return LoginRequiredAppError
	}

	// TODO: Quitar esto (ver #69)
	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		logger.Critical("UPLOAD_PATH no está definido")
	}

	articuloId := r.URL.Query().Get("articulo")
	if articuloId == "" {
		return &appError{ErrBadParam, "error obtaining `articulo`", "El parámetro de url `articulo` falta", 400}
	}

	files, err := os.ReadDir(uploadPath + "/extra")
	if err != nil {
		return &appError{err, "unable to read directory", "Error obteniendo datos", 500}
	}

	var result = make([]string, 0)

	for _, de := range files {
		if de.IsDir() {
			continue
		}

		name := de.Name()
		if strings.HasPrefix(name, articuloId) && strings.HasSuffix(name, ".jpg") {
			result = append(result, name)
		}
	}

	if len(result) < 1 {
		return &appError{errors.New("no matching pictures found"), "returning empty list", "No se ha encontrado ningún resultado", 404}
	}

	resbytes, err := json.MarshalIndent(result, "", "\t")
	if err != nil {
		return &appError{err, "error marhsalling response", "Error obteniendo datos", 500}
	}
	w.Write(resbytes)

	return nil
}

func subirFotoExtra(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		return LoginRequiredAppError
	}
	_, err = getSession(sc.Value)
	if err != nil {
		return LoginRequiredAppError
	}

	// TODO: Quitar esto (ver #69)
	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		logger.Critical("UPLOAD_PATH no está definido")
	}

	articuloId := r.FormValue("articulo")
	if articuloId == "" {
		return &appError{ErrBadParam, "error obtaining `articulo`", "El parámetro de url `articulo` falta", 400}
	}

	file, _, err := r.FormFile("foto")
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		return &appError{err, "error extracting image from form", "Hubo un error modificando la imagen.", 400}
	}

	photoBytes, err := io.ReadAll(file)
	if err != nil {
		return &appError{err, "error extracting image from form", "Hubo un error modificando la imagen.", 500}
	}

	image, err := imagenDesdeMime(photoBytes)
	if err != nil {
		return &appError{err, "error extracting image MIME", "Hubo un error modificando la imagen.", 500}
	}

	var fotoEscribir bytes.Buffer
	err = jpeg.Encode(&fotoEscribir, image, &jpeg.Options{Quality: 95})
	if err != nil {
		return &appError{err, "error encoding image", "Hubo un error modificando la imagen.", 500}
	}

	var salt = randstr.String(5)
	var imagePath = fmt.Sprintf("%s/extra/%s-%s.jpg", uploadPath, articuloId, salt)
	err = os.WriteFile(imagePath, fotoEscribir.Bytes(), 0o644)
	if err != nil {
		return &appError{err, "error writing image", "Hubo un error modificando la imagen.", 500}
	}

	return nil
}

func eliminarFotoExtra(w http.ResponseWriter, r *http.Request) *appError {
	var sc, err = r.Cookie("sess")
	if err != nil {
		return LoginRequiredAppError
	}
	sess, err := getSession(sc.Value)
	if err != nil {
		return LoginRequiredAppError
	}

	// TODO: Quitar esto (ver #69)
	uploadPath := os.Getenv("UPLOAD_PATH")
	if uploadPath == "" {
		logger.Critical("UPLOAD_PATH no está definido")
	}

	fotoId := r.URL.Query().Get("foto")
	if fotoId == "" {
		return &appError{ErrBadParam, "error obtaining `foto`", "El parámetro de url `foto` falta", 400}
	}

	/*
		Comprobar que el usuario que intenta eliminar la fotografía sea el autor, y que la publicación no esté pública
	*/

	row, err := db.Query(`SELECT COALESCE(fecha_publicacion, ""), autor_id FROM publicaciones WHERE id = ?`)
	if err != nil {
		return &appError{err, "error obtaining post data", "Error encontrando la fotografía a borrar", 500}
	}

	var dbFechaPub, dbAutorId string
	row.Scan(&dbFechaPub, &dbAutorId)

	if dbFechaPub != "" {
		return &appError{errors.New("post already public"), "post already public", "No se pueden eliminar fotografías de una publicación pública", 400}
	}

	// TODO: Crear permiso para saltarse esta limitación
	if sess.Autor_id != dbAutorId {
		return &appError{errors.New("author doesn't match"), "post already public", "No tienes permiso para borrar esta fotografía", 400}
	}

	/*
		Comprobar que la fotografía existe, y si es así eliminarla
	*/
	if f, err := os.Stat(uploadPath + "/extra/" + fotoId); err == nil {
		e2 := os.Remove(uploadPath + "/extra/" + f.Name())
		if e2 != nil {
			return &appError{e2, "error deleting file", "Hubo un error borrando la fotografía especificada", 500}
		}
		w.Write([]byte("{ \"error\": false }"))
	} else if errors.Is(err, os.ErrNotExist) {
		return &appError{err, "non-existent file", "El archivo especificado no existe", 400}
	} else {
		return &appError{err, "error finding file", "Hubo un error borrando la fotografía especificada", 500}
	}

	return nil
}
