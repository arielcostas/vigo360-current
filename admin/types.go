package admin

type Sesion struct {
	Id     string
	Nombre string
	Rol    string
}

type Aviso struct {
	Fecha_creacion string
	Titulo         string
	Contenido      string
}

type DashboardPost struct {
	Id                string
	Titulo            string
	Resumen           string
	Fecha_publicacion string
	Autor_nombre      string
}

type PostEditar struct {
	Id          string
	Titulo      string
	Resumen     string
	Contenido   string
	Alt_portada string
	Publicado   bool
}
