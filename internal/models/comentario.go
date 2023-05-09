package models

type EstadoComentario int

const (
	ESTADO_PENDIENTE EstadoComentario = 1
	ESTADO_APROBADO  EstadoComentario = 2
	ESTADO_RECHAZADO EstadoComentario = 3
)

type Comentario struct {
	Id             string
	Publicacion_id string
	Padre_id       string

	Nombre         string
	Es_autor       bool
	Autor_original bool
	Contenido      string

	Fecha_creacion   string
	Fecha_moderacion string
	Estado           EstadoComentario
	Moderador        string
}
