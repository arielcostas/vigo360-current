package service

import (
	"vigo360.es/new/internal/models"
)

type ComentarioTree struct {
	models.Comentario
}

func (se *Comentario) ListarPublicos(articulo_id string) ([]ComentarioTree, error) {
	// Obtiene de la base de datos
	comentariosLinear, err := se.cstore.ListarPublicos(articulo_id)
	if err != nil {
		return nil, err
	}

	// Convierte models.Comentario a ComentarioTree (esencialmente iguales)
	var clt = make([]ComentarioTree, 0)
	for _, c := range comentariosLinear {
		clt = append(clt, ComentarioTree{Comentario: c})
	}

	return clt, nil
}
