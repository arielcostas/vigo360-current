package admin

import (
	"net/http"

	"github.com/gorilla/mux"
)

func EditPostPage(w http.ResponseWriter, r *http.Request) {
	// TODO Check author is same as session
	verifyLogin(w, r)
	post_id := mux.Vars(r)["id"]
	post := PostEditar{}

	err := db.QueryRowx(`SELECT titulo, resumen, contenido, alt_portada FROM publicaciones WHERE id = ?;`, post_id).StructScan(&post)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error buscando el artículo en la base de datos"))
	}

	t.ExecuteTemplate(w, "post-id.html", struct {
		Post PostEditar
	}{Post: post})
}
