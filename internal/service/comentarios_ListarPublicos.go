// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package service

import (
	"vigo360.es/new/internal/models"
)

type ComentarioTree struct {
	models.Comentario
	Children []ComentarioTree
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

	// Llama a generarArbol para introducir los hijos en padre.Children
	return se.generarArbol(clt), nil
}
