package service

func (se *Comentario) generarArbol(ct []ComentarioTree) []ComentarioTree {
	var comentariosTreeMapeadosPorId = make(map[string]ComentarioTree, 0)
	for _, c := range ct {
		comentariosTreeMapeadosPorId[c.Id] = c
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
	return sliceComentarios
}
