package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/thanhpk/randstr"
	"golang.org/x/crypto/bcrypt"
	"vigo360.es/new/internal/database"
	"vigo360.es/new/internal/logger"
	"vigo360.es/new/internal/messages"
)

func (s *Server) handleAdminLoginAction() http.HandlerFunc {
	var comprobarContraseña = func(password string, hash string) bool {
		err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
		if err == nil {
			return true
		}

		if errors.Is(err, bcrypt.ErrHashTooShort) {
			// TODO: Refactor esto
			fmt.Printf("<3>contraseña demasiado corta")
		}

		return false
	}

	type LoginRow struct {
		Id         string
		Nombre     string
		Contraseña string
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

		if err := r.ParseForm(); err != nil {
			logger.Error("error recuperando datos del formulario: %s", err.Error())
			s.handleError(w, 400, messages.ErrorFormulario)
			return
		}

		var db = database.GetDB()
		param_userid := r.PostFormValue("userid")
		param_password := r.PostFormValue("password")

		row := LoginRow{}

		if param_userid == "" || param_password == "" {
			logger.Error("falta usuario o contraseña")
			s.handleAdminLoginPage(param_userid)(w, r)
			return
		}

		// Fetch from database. If error is no user found, show the error document. If a different error is thrown, show 500
		err = db.QueryRowx("SELECT id, nombre, contraseña FROM autores WHERE id=?;", param_userid).StructScan(&row)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				logger.Error("ningún usuario coincide con '%s'", param_userid)
				s.handleAdminLoginPage(param_userid)(w, r)
			} else {
				logger.Error("error recuperando usuario: %s", err.Error())
				s.handleError(w, 500, messages.ErrorDatos)
			}
			return
		}

		pass := comprobarContraseña(param_password, row.Contraseña)

		if !pass {
			logger.Error("la contraseña introducida para '%s' es inválida", param_userid)
			s.handleAdminLoginPage(param_userid)(w, r)
		}

		token := randstr.String(20)

		if _, err := db.Exec("INSERT INTO sesiones VALUES (?, NOW(), false, ?)", token, param_userid); err != nil {
			logger.Error("error guardando nueva sesión: %s")
			s.handleError(w, 500, messages.ErrorDatos)
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "sess",
			Value:    token,
			Path:     "/",
			MaxAge:   60*60*24*365,
			Domain:   r.URL.Host,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
		})

		var next = "/admin/dashboard"
		if n := r.URL.Query().Get("next"); n != "" {
			unescapedNext, err := url.QueryUnescape(n)
			if err == nil {
				next = unescapedNext
			}
		}

		defer w.WriteHeader(303)
		defer w.Header().Add("Location", next)
	}
}
