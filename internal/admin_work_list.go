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
type ResumenTrabajo struct {
	Id                string
	Titulo            string
	Fecha_publicacion string
	Publicado         bool
	Autor_id          string
	Autor_nombre      string
}

func (s *Server) handleAdminListTrabajos() http.HandlerFunc {
	const MAX_LISTED_WORKS = 50
	type response struct {
		Works   []ResumenTrabajo
		Session models.Session
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		db := database.GetDB()
		trabajos := []ResumenTrabajo{}

		err := db.Select(&trabajos, `SELECT trabajos.id,
       titulo,
       (fecha_publicacion < NOW() && fecha_publicacion IS NOT NULL) as publicado,
       COALESCE(fecha_publicacion, "")                              as fecha_publicacion,
       autor_id,
       autores.nombre                                               as autor_nombre
FROM trabajos
         LEFT JOIN autores ON trabajos.autor_id = autores.id
GROUP BY trabajos.id
ORDER BY publicado, trabajos.fecha_publicacion DESC;`)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			log.Error("error recuperando listado de trabajos: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorDatos)
			return
		}

		for i, p := range trabajos {
			tiempo, _ := time.Parse("2006-01-02 15:04:05", p.Fecha_publicacion)
			p.Fecha_publicacion = tiempo.Format("02/01/2006")
			trabajos[i] = p
		}

		if len(trabajos) > MAX_LISTED_WORKS {
			trabajos = trabajos[:MAX_LISTED_WORKS]
		}

		err = templates.Render(w, "admin-works.html", &response{
			Works:   trabajos,
			Session: sess,
		})

		if err != nil {
			log.Error("error recuperando el autor: %s", err.Error())
			s.handleError(r, w, 500, messages.ErrorRender)
		}
	}
}
