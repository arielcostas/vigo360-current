package internal

import (
	"net/http"

	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
	"vigo360.es/new/internal/templates"
)

func (s *Server) handleAdminLoginPage(prefill string) http.HandlerFunc {
	type response struct {
		LoginError  bool
		PrefillName string
	}

	return func(w http.ResponseWriter, r *http.Request) {
		logger := logger.NewLogger(r.Context().Value(ridContextKey("rid")).(string))
		var sc, err = r.Cookie("sess")
		if err == nil {
			sess, err := s.getSession(sc.Value)

			if err == nil { // User is logged in
				logger.Notice("%s ya tiene la sesión iniciada", sess.Autor_id)
				http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
				return
			}
		}

		var resp response
		if prefill != "" {
			resp = response{
				LoginError:  true,
				PrefillName: prefill,
			}
		} else {
			resp = response{}
		}

		err = templates.Render(w, "admin-login.html", resp)
		if err != nil {
			logger.Notice("error mostrando página: %s", err.Error())
			s.handleError(w, 500, messages.ErrorRender)
		}
	}
}
