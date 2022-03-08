package admin

import (
	"log"
	"net/http"

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
	err := t.ExecuteTemplate(w, "admin-login.html", &AdminLoginParams{})
	if err != nil {
		w.WriteHeader(500)
		log.Println("error with admin page: " + err.Error())
	}
}

func LoginAction(w http.ResponseWriter, r *http.Request) {
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
		//TODO log failed login
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

	// TODO Error handling
	db.Exec("INSERT INTO sesiones VALUES (?, NOW(), false, ?)", token, param_userid)

	http.SetCookie(w, &http.Cookie{
		Name:     "sess",
		Value:    token,
		Path:     "/admin",
		Domain:   r.URL.Host,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	println(param_userid + " logged in")
	http.Redirect(w, r, "/admin/dashboard", http.StatusTemporaryRedirect)
}
