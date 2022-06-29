package service

import (
	"errors"

	"github.com/thanhpk/randstr"
	"vigo360.es/new/internal/models"
)

var Err_ComentarioPublicacionInvalida = errors.New("el ID de la publicación no es válido")
var Err_ComentarioNombreInvalido = errors.New("el nombre de autor del comentario no es válido")
var Err_ComentarioContenidoInvalido = errors.New("el contenido del comentario no es válido")

func (se *Comentario) AgregarComentario(
	publicacion_id string,
	nombre string,
	contenido string,
) (models.Comentario, error) {
	if nombre == "" || len(nombre) > 40 {
		return models.Comentario{}, Err_ComentarioNombreInvalido
	}

	if contenido == "" || len(contenido) > 500 {
		return models.Comentario{}, Err_ComentarioContenidoInvalido
	}

	publicacion_existe, err := se.pstore.Existe(publicacion_id)
	if err != nil {
		return models.Comentario{}, err
	}
	if !publicacion_existe {
		return models.Comentario{}, Err_ComentarioPublicacionInvalida
	}

	var nuevo_comentario = models.Comentario{
		Id:             randstr.String(13),
		Publicacion_id: publicacion_id,

		Nombre:         nombre,
		Es_autor:       false,
		Autor_original: false,
		Contenido:      contenido,

		Estado: models.ESTADO_PENDIENTE,
	}

	return nuevo_comentario, se.cstore.GuardarComentario(nuevo_comentario)
}

// func (se *Comentario) AgregarRespuesta(
// 	publicacion_id string,
// 	nombre string,
// 	contenido string,
// 	padre_id string,
// ) (models.Comentario, error) {

// }
