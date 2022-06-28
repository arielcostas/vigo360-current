package service

import (
	"vigo360.es/new/internal/models"
)

type ComentarioTree struct {
	models.Comentario
	Children []ComentarioTree
}

func (se *Comentario) ListarPublicos(articulo_id string) ([]ComentarioTree, error) {
	comentariosLinear, err := se.store.ListarPublicos(articulo_id)
	if err != nil {
		return nil, err
	}

	var comentariosTreeMapeadosPorId = make(map[string]ComentarioTree, 0)
	for _, c := range comentariosLinear {
		comentariosTreeMapeadosPorId[c.Id] = ComentarioTree{Comentario: c}
	}

	var mapaComentarios = make(map[string]ComentarioTree)

	for _, c := range comentariosTreeMapeadosPorId {
		if c.Padre_id == "" {
			mapaComentarios[c.Id] = c
		} else {
			var padre = comentariosTreeMapeadosPorId[c.Padre_id]
			padre.Children = append(padre.Children, c)
			mapaComentarios[padre.Id] = padre
		}
	}

	var sliceComentarios = make([]ComentarioTree, 0)
	for _, ct := range mapaComentarios {
		sliceComentarios = append(sliceComentarios, ct)
	}

	return sliceComentarios, nil
}
