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

func GetFullTrabajo(id string) (Trabajo, error) {
	query := `SELECT trabajos.id, alt_portada, titulo, resumen, contenido, 
	DATE_FORMAT(trabajos.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	DATE_FORMAT(trabajos.fecha_publicacion, '%d %b. %Y') as fecha_actualizacion, 
	autores.id as autor_id, autores.nombre as autor_nombre, autores.biografia as autor_biografia, autores.rol as autor_rol
	FROM TrabajosPublicos trabajos 
	LEFT JOIN autores on trabajos.autor_id = autores.id 
	WHERE trabajos.id = ?;`

	trabajo := Trabajo{}
	if err := db.QueryRowx(query, id).StructScan(&trabajo); err != nil {
		return Trabajo{}, err
	}

	return trabajo, nil
}

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

	rows, err := db.Query(query)
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
