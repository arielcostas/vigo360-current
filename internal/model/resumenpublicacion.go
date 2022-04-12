/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package model

import "vigo360.es/new/internal/database"

type ResumenPublicacion struct {
	Id                  string
	Fecha_publicacion   string
	Fecha_actualizacion string
	Alt_portada         string
	Titulo              string
	Resumen             string
	Autor               struct {
		Id     string
		Nombre string
		Email  string
	}
	Tags string
}

func ListarPublicacionesPublicas() ([]ResumenPublicacion, error) {
	rp := make([]ResumenPublicacion, 0)
	query := `SELECT pp.id, fecha_publicacion, fecha_actualizacion, titulo, resumen, autor_id, autores.nombre as autor_nombre, autores.email as autor_email, GROUP_CONCAT(tags.nombre) as tags FROM PublicacionesPublicas pp LEFT JOIN publicaciones_tags ON pp.id = publicaciones_tags.publicacion_id LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id LEFT JOIN autores ON pp.autor_id = autores.id GROUP BY id ORDER BY fecha_publicacion;`

	rows, err := database.GetDB().Query(query)
	if err != nil {
		return rp, err
	}
	for rows.Next() {
		var np ResumenPublicacion
		rows.Scan(&np.Id, &np.Fecha_publicacion, &np.Fecha_actualizacion, &np.Titulo, &np.Resumen, &np.Autor.Id, &np.Autor.Nombre, &np.Autor.Email, &np.Tags)
		rp = append(rp, np)
	}

	return rp, nil
}
