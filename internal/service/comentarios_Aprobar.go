package service

func (se *Comentario) Aprobar(comentario_id string, moderador_id string) error {
	return se.cstore.Aprobar(comentario_id, moderador_id)
}
