package repository

import (
	"crud/internal/models"
	"database/sql"
)

type UsuarioRepository interface {
	ObtenerTodos() ([]models.Usuario, error)
	Crear(u models.Usuario) (models.Usuario, error)
	Actualizar(u models.Usuario) error
	Eliminar(id int) error
}

type mysqlUsuarioRepository struct {
	db *sql.DB
}

func NewMySQLUsuarioRepository(db *sql.DB) UsuarioRepository {
	return &mysqlUsuarioRepository{db: db}
}

func (r *mysqlUsuarioRepository) ObtenerTodos() ([]models.Usuario, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM usuarios")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var u models.Usuario
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	if usuarios == nil {
		usuarios = []models.Usuario{}
	}
	return usuarios, nil
}

func (r *mysqlUsuarioRepository) Crear(u models.Usuario) (models.Usuario, error) {
	res, err := r.db.Exec("INSERT INTO usuarios (name, email) VALUES (?, ?)", u.Name, u.Email)
	if err != nil {
		return u, err
	}
	id, _ := res.LastInsertId()
	u.ID = int(id)
	return u, nil
}

func (r *mysqlUsuarioRepository) Actualizar(u models.Usuario) error {
	_, err := r.db.Exec("UPDATE usuarios SET name = ?, email = ? WHERE id = ?", u.Name, u.Email, u.ID)
	return err
}

func (r *mysqlUsuarioRepository) Eliminar(id int) error {
	_, err := r.db.Exec("DELETE FROM usuarios WHERE id = ?", id)
	return err
}
