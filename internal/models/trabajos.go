/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package models

import "time"

type Trabajos []Trabajo

// Devuelve un slice con solo los trabajos públicos
func (ps *Trabajos) FiltrarPublicos() Trabajos {
	var nps Trabajos

	for _, p := range *ps {
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

// Devuelve la fecha de la actualización más reciente de una publicación del slice
func (ps *Trabajos) ObtenerUltimaActualizacion() (time.Time, error) {
	var lastUpdate time.Time
	for _, pub := range *ps {
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
