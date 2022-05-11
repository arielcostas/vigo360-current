package internal

import (
	"database/sql"
	"errors"
	"net/http"

	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/models"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handleAdminDashboardPage() http.HandlerFunc {
	type response struct {
		Avisos  []Aviso
		Posts   []DashboardPost
		Session models.Session
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		sess := r.Context().Value(sessionContextKey("sess")).(models.Session)
		db := database.GetDB()

		avisos := []Aviso{}
		err := db.Select(&avisos, "SELECT DATE_FORMAT(fecha_creacion, '%d %b.') as fecha_creacion, titulo, contenido FROM avisos ORDER BY avisos.fecha_creacion DESC LIMIT 5")

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando avisos: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
		}

		posts := []DashboardPost{}
		err = db.Select(&posts, "SELECT publicaciones.id, titulo, DATE_FORMAT(fecha_publicacion, '%d %b.') as fecha_publicacion, resumen, autores.nombre as autor_nombre FROM publicaciones LEFT JOIN autores ON publicaciones.autor_id = autores.id WHERE publicaciones.fecha_publicacion IS NOT NULL ORDER BY publicaciones.fecha_publicacion DESC LIMIT 5;")

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			logger.Error("error recuperando últimas publicaciones: %s", err.Error())
			s.handleError(w, 500, messages.ErrorDatos)
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
