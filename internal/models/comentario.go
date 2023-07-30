package models

type EstadoComentario int

//goland:noinspection ALL
const (
	EstadoPendiente EstadoComentario = 1
	EstadoAprobado  EstadoComentario = 2
	EstadoRechazado EstadoComentario = 3
)

type Comentario struct {
	Id                 string
	Publicacion_id     string
	Publicacion_titulo string
	Padre_id           string

	Nombre         string
	Es_autor       bool
	Autor_original bool
	Contenido      string

	Fecha_creacion   string
	Fecha_moderacion string
	Estado           EstadoComentario
	Moderador        string
}
