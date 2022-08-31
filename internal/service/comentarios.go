// SPDX-FileCopyrightText: 2022 Ariel Costas <ariel@vigo360.es>
//
// SPDX-License-Identifier: MPL-2.0

package service

import "vigo360.es/new/internal/repository"

type Comentario struct {
	cstore repository.ComentarioStore
	pstore repository.PublicacionStore
}

func NewComentarioService(cstore repository.ComentarioStore, pstore repository.PublicacionStore) Comentario {
	return Comentario{cstore: cstore, pstore: pstore}
}
