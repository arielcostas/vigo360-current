/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

import "time"

type Publicaciones []Publicacion

func (ps *Publicaciones) FiltrarPublicas() Publicaciones {
	var nps Publicaciones

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
