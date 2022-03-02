package routes

import (
	"bytes"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

//#region
type TrabajosTrabajo struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string
	Alt_portada         string
	Titulo              string
	Resumen             string
	ContenidoRaw        string `db:"contenido"`
	Contenido           template.HTML
	Autor_id            string
	Autor_nombre        string
	Autor_rol           string
	Autor_biografia     string
}

type TrabajoAdjunto struct {
	Nombre_archivo string
	Titulo         string
}

type TrabajoParams struct {
	Trabajo  TrabajosTrabajo
	Adjuntos []TrabajoAdjunto
	Meta     PageMeta
}

func TrabajoPage(w http.ResponseWriter, r *http.Request) {
	req_paper_id := mux.Vars(r)["paperid"]
	query := `SELECT trabajos.id, alt_portada, titulo, resumen, contenido, 
	DATE_FORMAT(trabajos.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	DATE_FORMAT(trabajos.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol
FROM trabajos 
LEFT JOIN autores on trabajos.autor_id = autores.id 
WHERE trabajos.id = ?;`

	trabajo := TrabajosTrabajo{}
	row := db.QueryRowx(query, req_paper_id)
	err := row.StructScan(&trabajo)
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Result is in markdown, convert to HTML
	var buf bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(extension.Footnote))
	err = md.Convert([]byte(trabajo.ContenidoRaw), &buf)
	if err != nil {
		log.Fatalf(err.Error())
	}
	trabajo.Contenido = template.HTML(buf.Bytes())

	adjuntos := []TrabajoAdjunto{}
	err = db.Select(&adjuntos, "SELECT nombre_archivo, titulo FROM adjuntos WHERE trabajo_id = ?;", trabajo.Id)

	if err != nil {
		log.Fatalf(err.Error())
	}

	t.ExecuteTemplate(w, "trabajos-id.html", TrabajoParams{
		Trabajo:  trabajo,
		Adjuntos: adjuntos,
		Meta: PageMeta{
			Title: trabajo.Titulo,
		},
	})
}

//#endregion

type TrabajoListParams struct {
	Trabajos []TrabajosListTrabajo
	Meta     PageMeta
}

type TrabajosListTrabajo struct {
	Id                string
	Fecha_publicacion string
	Alt_portada       string
	Titulo            string
	Autor_nombre      string
}

func TrabajoListPage(w http.ResponseWriter, r *http.Request) {
	trabajos := []TrabajosListTrabajo{}
	err := db.Select(&trabajos, `SELECT trabajos.id, DATE_FORMAT(fecha_publicacion, '%d %b. %Y') as fecha_publicacion, alt_portada, titulo, autores.nombre as autor_nombre FROM trabajos LEFT JOIN autores on trabajos.autor_id = autores.id WHERE trabajos.fecha_publicacion < NOW() ORDER BY trabajos.fecha_publicacion DESC`)

	if err != nil {
		log.Fatalf(err.Error())
	}

	t.ExecuteTemplate(w, "trabajos.html", TrabajoListParams{
		Trabajos: trabajos,
		Meta: PageMeta{
			Title: "Trabajos",
		},
	})
}
