package service

import "vigo360.es/new/internal/repository"

type Comentario struct {
	store repository.ComentarioStore
}

func NewComentarioService(store repository.ComentarioStore) Comentario {
	return Comentario{store: store}
}
