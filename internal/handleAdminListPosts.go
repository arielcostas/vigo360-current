package internal

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

// TODO: Get rid of this
type ResumenPost struct {
	Id                string
	Titulo            string
	Fecha_publicacion string
	CantTags          int
	Publicado         bool
	Autor_id          string
	Autor_nombre      string
}

func (s *Server) handleAdminListPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		db := database.GetDB()
		posts := []ResumenPost{}

		err := db.Select(&posts, `SELECT publicaciones.id, titulo, (fecha_publicacion < NOW() && fecha_publicacion IS NOT NULL) as publicado, COALESCE(fecha_publicacion, "") as fecha_publicacion, autor_id, autores.nombre as autor_nombre, count(tag_id) as canttags FROM publicaciones LEFT JOIN autores ON publicaciones.autor_id = autores.id LEFT JOIN publicaciones_tags ON publicaciones.id = publicaciones_tags.publicacion_id GROUP BY publicaciones.id ORDER BY publicado ASC, publicaciones.fecha_publicacion DESC;`)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando listado de publicacionoes: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
			return
		}

		for i, p := range posts {
			tiempo, _ := time.Parse("2006-01-02 15:04:05", p.Fecha_publicacion)
			p.Fecha_publicacion = tiempo.Format("02/01/2006")
			posts[i] = p
		}

		err = templates.Render(w, "admin-post.html", struct {
			Posts   []ResumenPost
			Session models.Session
		}{
			Posts:   posts,
			Session: sess,
		})

		if err != nil {
			logger.Error("error recuperando el autor: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
