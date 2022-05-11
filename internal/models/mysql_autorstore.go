package models

import "github.com/jmoiron/sqlx"

type MysqlAutorStore struct {
	db *sqlx.DB
}

func NewMysqlAutorStore(db *sqlx.DB) *MysqlAutorStore {
	return &MysqlAutorStore{
		db: db,
	}
}

func (s *MysqlAutorStore) Listar() ([]Autor, error) {
	var autores = make([]Autor, 0)
	var rows, err = s.db.Query(`SELECT id, nombre, email, rol, biografia FROM autores`)
	if err != nil {
		return []Autor{}, err
	}

	for rows.Next() {
		var na Autor
		err = rows.Scan(&na.Id, &na.Nombre, &na.Email, &na.Rol, &na.Biografia)
		if err != nil {
			return []Autor{}, err
		}
		autores = append(autores, na)
	}

	return autores, nil
}

func (s *MysqlAutorStore) Obtener(autor_id string) (Autor, error) {
	var autor Autor
	var row = s.db.QueryRow(`SELECT id, nombre, email, rol, biografia, web_url, web_titulo FROM autores WHERE id=?`, autor_id)
	var err = row.Scan(&autor.Id, &autor.Nombre, &autor.Email, &autor.Rol, &autor.Biografia, &autor.Web.Url, &autor.Web.Titulo)

	if err != nil {
		return Autor{}, err
	}

	return autor, nil
}

func (s *MysqlAutorStore) Buscar(termino string) ([]Autor, error) {
	var autores []Autor

	var query = `SELECT id, nombre, email, rol, biografia FROM autores WHERE CONCAT(nombre, email, rol, biografia) LIKE ?`
	var rows, err = s.db.Query(query, "%"+termino+"%")
	if err != nil {
		return []Autor{}, err
	}

	for rows.Next() {
		var na Autor
		err = rows.Scan(&na.Id, &na.Nombre, &na.Email, &na.Rol, &na.Biografia)
		if err != nil {
			return []Autor{}, err
		}
		autores = append(autores, na)
	}

	return autores, nil
}
