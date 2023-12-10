package models

import (
	"time"
)

type Publicaciones []Publicacion

// Devuelve un slice con solo las publicaciones públicas
func (ps Publicaciones) FiltrarPublicas() Publicaciones {
	var nps Publicaciones

	for _, p := range ps {
		if p.Fecha_publicacion == "" {
			continue
		}

		var fechaPub, err = time.Parse("2006-01-02 15:04:05" /* Y-M-D H:M:S*/, p.Fecha_publicacion)
		if err != nil {
			continue
		}
		if fechaPub.Unix() <= time.Now().Unix() {
			nps = append(nps, p)
		}
	}

	return nps
}

// FiltrarRetiradas Devuelve un slice con solo las publicaciones retiradas por razones legales
func (ps Publicaciones) FiltrarRetiradas() Publicaciones {
	var nps Publicaciones

	for _, p := range ps {
		if p.Legally_retired_at != "" {
			continue
		}

		nps = append(nps, p)
	}

	return nps
}

// Devuelve la fecha de la actualización más reciente de una publicación del slice
func (ps Publicaciones) ObtenerUltimaActualizacion() (time.Time, error) {
	var lastUpdate time.Time
	for _, pub := range ps {
		var ut, err = time.Parse("2006-01-02 15:04:05", pub.Fecha_actualizacion)
		if err != nil {
			return time.Time{}, err
		}

		if ut.Unix() > lastUpdate.Unix() {
			lastUpdate = ut
		}
	}

	return lastUpdate, nil
}
