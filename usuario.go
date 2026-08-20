package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// ============================================================================
// 1. MODELO / ENTIDAD
// ============================================================================

type Usuario struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ============================================================================
// 2. CAPA DE DATOS (REPOSITORIO) - Abstracción con Interfaz
// ============================================================================

// UsuarioRepository define el contrato (Polimorfismo / Abstracción)
type UsuarioRepository interface {
	ObtenerTodos() ([]Usuario, error)
	Crear(u Usuario) (Usuario, error)
	Actualizar(u Usuario) error
	Eliminar(id int) error
}

// mysqlUsuarioRepository es la implementación concreta para MySQL (Encapsulamiento)
type mysqlUsuarioRepository struct {
	db *sql.DB
}

// Constructor del Repositorio
func NewMySQLUsuarioRepository(db *sql.DB) UsuarioRepository {
	return mysqlUsuarioRepository{db: db}
}

func (r mysqlUsuarioRepository) ObtenerTodos() ([]Usuario, error) {
	var rows, err = r.db.Query("SELECT id, name, email FROM usuarios")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var u Usuario
		var err = rows.Scan(&u.ID, &u.Name, &u.Email)
		if err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}

	var err2 = rows.Err()
	if err2 != nil {
		return nil, err2
	}

	if usuarios == nil {
		usuarios = []Usuario{}
	}

	return usuarios, nil
}

func (r mysqlUsuarioRepository) Crear(u Usuario) (Usuario, error) {
	var res, errE = r.db.Exec("INSERT INTO usuarios (name, email) VALUES (?, ?)", u.Name, u.Email)
	if errE != nil {
		return u, errE
	}
	var id, err = res.LastInsertId()
	if err != nil {
		return u, err
	}
	u.ID = int(id)
	return u, nil
}

func (r mysqlUsuarioRepository) Actualizar(u Usuario) error {
	var _, err = r.db.Exec("UPDATE usuarios SET name = ?, email = ? WHERE id = ?", u.Name, u.Email, u.ID)
	return err
}

func (r mysqlUsuarioRepository) Eliminar(id int) error {
	var _, err = r.db.Exec("DELETE FROM usuarios WHERE id = ?", id)
	return err
}

// ============================================================================
// 3. CAPA DE PRESENTACIÓN (CONTROLADOR / HANDLER HTTP)
// ============================================================================

// UsuarioHandler maneja únicamente las peticiones HTTP (Inyección de Dependencias)
type UsuarioHandler struct {
	repo UsuarioRepository // Depende de la interfaz, NO de MySQL directamente
}

// Constructor del Handler
func NewUsuarioHandler(repo UsuarioRepository) UsuarioHandler {
	return UsuarioHandler{repo: repo}
}

func (h UsuarioHandler) HandleUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		var usuarios, err = h.repo.ObtenerTodos()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(usuarios)

	case "POST":
		var u Usuario
		var err = json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var usuarioCreado, err2 = h.repo.Crear(u)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(usuarioCreado)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

func (h UsuarioHandler) HandleUsuarioPorID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var idStr = r.URL.Path[len("/api/usuarios/"):]
	var id, err = strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "PUT":
		var u Usuario
		var err = json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.ID = id
		var err2 = h.repo.Actualizar(u)
		if err2 != nil {
			http.Error(w, err2.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(u)

	case "DELETE":
		var err = h.repo.Eliminar(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}
