package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// ==========================================
// 1. ENTIDAD (Objeto / Modelo)
// ==========================================

// Usuario representa nuestra entidad en la BDD
type Usuario struct {
	ID     int
	Nombre string
	Email  string
}

// ==========================================
// 2. INTERFAZ (Contrato de comportamiento)
// ==========================================

type UsuarioRepository interface {
	Crear(u *Usuario) error
	ObtenerPorID(id int) (*Usuario, error)
	ObtenerTodos() ([]Usuario, error)
	Actualizar(u *Usuario) error
	Eliminar(id int) error
}

// ==========================================
// 3. IMPLEMENTACIÓN POO (Estructura con métodos)
// ==========================================

type mysqlUsuarioRepository struct {
	db *sql.DB // Encapsulamos la conexión
}

// Constructor para instanciar el repositorio
func NewUsuarioRepository(db *sql.DB) UsuarioRepository {
	return &mysqlUsuarioRepository{db: db}
}

// CREATE
func (r *mysqlUsuarioRepository) Crear(u *Usuario) error {
	query := "INSERT INTO usuarios (nombre, email) VALUES (?, ?)"
	result, err := r.db.Exec(query, u.Nombre, u.Email)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(id)
	return nil
}

// READ (Uno)
func (r *mysqlUsuarioRepository) ObtenerPorID(id int) (*Usuario, error) {
	query := "SELECT id, nombre, email FROM usuarios WHERE id = ?"
	row := r.db.QueryRow(query, id)

	u := &Usuario{}
	err := row.Scan(&u.ID, &u.Nombre, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, err
	}
	return u, nil
}

// READ (Todos)
func (r *mysqlUsuarioRepository) ObtenerTodos() ([]Usuario, error) {
	query := "SELECT id, nombre, email FROM usuarios"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var u Usuario
		if err := rows.Scan(&u.ID, &u.Nombre, &u.Email); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, u)
	}
	return usuarios, nil
}

// UPDATE
func (r *mysqlUsuarioRepository) Actualizar(u *Usuario) error {
	query := "UPDATE usuarios SET nombre = ?, email = ? WHERE id = ?"
	_, err := r.db.Exec(query, u.Nombre, u.Email, u.ID)
	return err
}

// DELETE
func (r *mysqlUsuarioRepository) Eliminar(id int) error {
	query := "DELETE FROM usuarios WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

// ==========================================
// 4. EJECUCIÓN DEL PROGRAMA (Main)
// ==========================================

func main() {
	// Cadena de conexión: usuario:contraseña@tcp(127.0.0.1:3306)/nombre_bdd
	dsn := "root:@tcp(127.0.0.1:3306)/empresa"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error de configuración de BDD: %v", err)
	}
	defer db.Close()

	// Probar conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo conectar a la BDD: %v", err)
	}
	fmt.Println("¡Conexión exitosa a MySQL!")

	// Instanciamos el repositorio (Inyección de dependencias)
	repo := NewUsuarioRepository(db)

	// --- 1. CREATE ---
	nuevoUsuario := &Usuario{Nombre: "Carlos Pérez", Email: "carlos@example.com"}
	err = repo.Crear(nuevoUsuario)
	if err != nil {
		log.Printf("Error al crear: %v", err)
	} else {
		fmt.Printf("Usuario creado con ID: %d\n", nuevoUsuario.ID)
	}

	// --- 2. READ (Por ID) ---
	usuario, err := repo.ObtenerPorID(nuevoUsuario.ID)
	if err != nil {
		log.Printf("Error al obtener: %v", err)
	} else {
		fmt.Printf("Usuario obtenido: %+v\n", usuario)
	}

	// --- 3. UPDATE ---
	usuario.Nombre = "Carlos A. Pérez"
	err = repo.Actualizar(usuario)
	if err != nil {
		log.Printf("Error al actualizar: %v", err)
	} else {
		fmt.Println("Usuario actualizado con éxito")
	}

	// --- 2. READ ALL ---
	todos, err := repo.ObtenerTodos()
	if err == nil {
		fmt.Println("\nLista de todos los usuarios:")
		for _, u := range todos {
			fmt.Printf("- ID: %d, Nombre: %s, Email: %s\n", u.ID, u.Nombre, u.Email)
		}
	}

	// --- 4. DELETE ---
	err = repo.Eliminar(usuario.ID)
	if err != nil {
		log.Printf("Error al eliminar: %v", err)
	} else {
		fmt.Printf("Usuario con ID %d eliminado con éxito\n", usuario.ID)
	}
}
