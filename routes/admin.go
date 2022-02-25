package routes

import (
	"fmt"
	"net/http"

	_ "embed"
)

//go:embed templates/admin/login.html
var adminLoginPage []byte

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	// Serve the form and end
	if r.Method == http.MethodGet {
		w.Write(adminLoginPage)
		return
	}

	r.ParseForm()

	var row struct {
		id       string
		nombre   string
		password string
	}
	err := db.QueryRow("SELECT id, nombre, contraseña FROM autores WHERE id=?;", r.PostFormValue("userid")).Scan(&row.id, &row.nombre, &row.password)
	if err != nil {
		println(err.Error())
	}

	pass := ValidPassword(r.PostFormValue("password"), row.password)
	fmt.Fprintf(w, `id => %s
nombre => %s
contraseña => %t
	`, row.id, row.nombre, pass)

}
