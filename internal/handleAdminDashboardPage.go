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

func (s *Server) handleAdminDashboardPage() http.HandlerFunc {
	type response struct {
		Avisos  []models.Aviso
		Posts   []models.Publicacion
		Session models.Session
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		db := database.GetDB()

		// TODO: Convertir esto en llamada a repositorio
		avisos := []models.Aviso{}
		err := db.Select(&avisos, "SELECT fecha_creacion, titulo, contenido FROM avisos ORDER BY fecha_creacion DESC LIMIT 5")

		for i, a := range avisos {
			tiempo, _ := time.Parse("2006-01-02 15:04:05", a.Fecha_creacion)
			a.Fecha_creacion = tiempo.Format("02/01")
			avisos[i] = a
		}

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando avisos: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
		}

		var posts models.Publicaciones
		posts, err = s.store.publicacion.Listar()

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando últimas publicaciones: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
		}

		posts = posts.FiltrarPublicas()[0:4]

		for i, p := range posts {
			tiempo, _ := time.Parse("2006-01-02 15:04:05", p.Fecha_publicacion)
			p.Fecha_publicacion = tiempo.Format("02/01")
			posts[i] = p
		}

		err = templates.Render(w, "admin-dashboard.html", response{
			Avisos:  avisos,
			Posts:   posts,
			Session: sess,
		})
		if err != nil {
			logger.Error("error mostrando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
