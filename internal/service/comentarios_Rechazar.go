package service

func (se *Comentario) Rechazar(comentario_id string, moderador_id string) error {
	return se.cstore.Rechazar(comentario_id, moderador_id)
}
