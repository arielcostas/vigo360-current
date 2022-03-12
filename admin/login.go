package admin

import (
	"net/http"

	"git.sr.ht/~arielcostas/new.vigo360.es/logger"
	"github.com/thanhpk/randstr"
)

type AdminLoginParams struct {
	LoginError  bool
	PrefillName string
}

type LoginRow struct {
	Id         string
	Nombre     string
	Contraseña string
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)
	err := t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{})
	if err != nil {
		logger.Error("[adminlogin]: error rendering page: %s", err.Error())
		w.WriteHeader(500)
		InternalServerErrorHandler(w, r)
	}
}

func LoginAction(w http.ResponseWriter, r *http.Request) {
	verifyLogin(w, r)

	r.ParseForm()
	param_userid := r.PostFormValue("userid")
	param_password := r.PostFormValue("password")

	row := LoginRow{}

	if param_userid == "" || param_password == "" {
		t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{
			PrefillName: param_userid,
			LoginError:  true,
		})
		return
	}

	err := db.QueryRowx("SELECT id, nombre, contraseña FROM autores WHERE id=?;", param_userid).StructScan(&row)

	if err != nil {
		logger.Error("[login] failed login for user %s", param_userid)
		t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{
			PrefillName: param_userid,
			LoginError:  true,
		})
		return
	}

	pass := ValidatePassword(param_password, row.Contraseña)

	if !pass {
		t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{
			PrefillName: param_userid,
			LoginError:  true,
		})
		return
	}

	token := randstr.String(20)

	_, err = db.Exec("INSERT INTO sesiones VALUES (?, NOW(), false, ?)", token, param_userid)

	if err != nil {
		logger.Error("[login] error saving new session token %s for user %s", token, param_userid)
		InternalServerErrorHandler(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    token,
		Path:     "/admin",
		Domain:   r.URL.Host,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	w.Header().Add("Location", "/admin/dashboard")
	w.WriteHeader(http.StatusSeeOther)
}
