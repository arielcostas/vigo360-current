package admin

import (
	"net/http"
)

func DashboardPage(w http.ResponseWriter, r *http.Request) {
	sesion := verifyLogin(w, r)

	t.ExecuteTemplate(w, "admin-dashboard.html", struct {
		Nombre string
	}{
		Nombre: sesion.Nombre,
	})
}
