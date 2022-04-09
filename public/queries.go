/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package public

func GetFullPost(id string) (FullPost, error) {
	var query = `SELECT 
	pp.id, alt_portada, titulo, resumen, contenido, DATE_FORMAT(pp.fecha_publicacion, '%d %b.') as fecha_publicacion, DATE_FORMAT(pp.fecha_actualizacion, '%d %b.') as fecha_actualizacion, autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol, serie_id as serie, GROUP_CONCAT(tags.nombre) as tags 
    FROM PublicacionesPublicas pp
    LEFT JOIN autores on pp.autor_id = autores.id
    LEFT JOIN publicaciones_tags ON pp.id = publicaciones_tags.publicacion_id
    LEFT JOIN tags ON publicaciones_tags.tag_id = tags.id
    WHERE pp.id = ?
    GROUP BY pp.id 
    ORDER BY pp.fecha_publicacion DESC;`

	post := FullPost{}
	if err := db.QueryRowx(query, id).StructScan(&post); err != nil {
		return FullPost{}, err
	} else {
		return post, nil
	}
}

func GetSerieById(id string) (Serie, error) {
	var serie Serie
	if err := db.QueryRowx(`SELECT titulo FROM series WHERE id = ?;`, id).Scan(&serie.Titulo); err != nil {
		return Serie{}, err
	}

	if err := db.Select(&serie.Articulos, `SELECT id, titulo FROM PublicacionesPublicas WHERE serie_id=? ORDER BY serie_posicion ASC, titulo ASC`, id); err != nil {
		return Serie{}, err
	}

	return serie, nil
}
