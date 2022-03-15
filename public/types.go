package public

import (
	"database/sql"
	"html/template"
)

type Autor struct {
	Id         string
	Nombre     string
	Email      string
	Rol        string
	Biografia  string
	Web_url    string
	Web_titulo string
}

type ResumenPost struct {
	Id                string
	Fecha_publicacion string
	Alt_portada       string
	Titulo            string
	Resumen           string
	Autor_id          string
	Autor_nombre      string `db:"nombre"`
	Serie_id          string
	Serie_posicion    int
}

type FullPost struct {
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
	Serie               sql.NullString
}

type Serie struct {
	Titulo    string
	Articulos []ResumenPost
}

type Tag struct {
	Id     string
	Nombre string
}
