package internal

import (
	"database/sql"
	_ "embed"
	"github.com/go-playground/validator/v10"
	"net/http"
	"os"
	"regexp"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
)

//go:embed extra/default.jpg
var defaultImageJPG []byte

//go:embed extra/default.webp
var defaultImageWebp []byte

func (s *Server) handleAdminCreatePost() http.HandlerFunc {
	type entrada struct {
		ArtId  string
		Titulo string `validate:"required,min=3,max=80"`
	}

	type response struct {
		Posts   []ResumenPost
		Session models.Session
	}

	var postIdRegexp = regexp.MustCompile(`^[A-Za-z0-9\-_]{3,40}$`)
	var ValidatePostId = func(id string) bool {
		return postIdRegexp.MatchString(id)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess, _ := r.Context().Value(sessionContextKey("sess")).(models.Session)

		var artAutor = sess.Autor_id

		if err := r.ParseForm(); err != nil {
			log.Error("error leyendo datos de formulario: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorFormulario)
			return
		}

		fi := entrada{
			ArtId:  r.FormValue("art-id"),
			Titulo: r.FormValue("art-titulo"),
		}

		if !ValidatePostId(fi.ArtId) {
			log.Error("error validando id")
			s.handleError(r, w, 500, messages.ErrorIdInvalido)
			return
		}

		if err := validator.New().Struct(fi); err != nil {
			// TODO: Mostrar de nuevo formulario con los errores
			log.Error("error validando título: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorValidacion)
			return
		}

		var tx *sql.Tx
		if nt, err := database.GetDB().Begin(); err != nil {
			log.Error("error iniciando transacción: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		} else {
			tx = nt
		}

		q := "INSERT INTO publicaciones(id, titulo, alt_portada, resumen, contenido, autor_id) VALUES (?, ?, 'CAMBIAME','', '', ?);"
		if _, err := tx.Exec(q, fi.ArtId, fi.Titulo, artAutor); err != nil {
			tx.Rollback()
			log.Error("error creando artículo: %s", err.Error())
			if err.Error() == "UNIQUE constraint failed: publicaciones.id" {
				s.handleError(r, w, 500, messages.ErrorIdDuplicado)
			} else {
				s.handleError(r, w, 500, messages.ErrorDatos)
			}
			return
		}

		photopath := os.Getenv("UPLOAD_PATH")
		if err := os.WriteFile(photopath+"/images/"+fi.ArtId+".webp", defaultImageWebp, 0o644); err != nil {
			tx.Rollback()
			log.Error("error escribiendo imagen webp: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		if err := os.WriteFile(photopath+"/thumb/"+fi.ArtId+".jpg", defaultImageJPG, 0o644); err != nil {
			tx.Rollback()
			log.Error("error escribiendo imagen jpg: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			log.Error("error ejecutando transacción: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		w.Header().Add("Location", "/admin/post/"+fi.ArtId)
		w.WriteHeader(303)
	}
}
