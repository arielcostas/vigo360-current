package admin

import "regexp"

func verificarId(id string) bool {
	return regexp.MustCompile(`^[a-z|ñ|\-|\_|0-9]{3,40}$`).MatchString(id)
}

func verificarTitulo(titulo string) bool {
	return len(titulo) > 3 && len(titulo) < 80
}

func verificarResumen(resumen string) bool {
	return len(resumen) > 3 && len(resumen) < 300
}
