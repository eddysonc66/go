package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// ==========================================
// 1. ENTIDAD (Objeto / Modelo)
// ==========================================

// Usuario representa nuestra entidad en la BDD
// Nota: Usuario: Define cómo se ve un usuario en memoria. Los fragmentos json:"id" se llaman Struct Tags y le dicen a Go qué nombre usar en las claves del objeto JSON que viajará al cliente web.
// Nota 2: Los campos de una estructura solo son visibles para otros paquetes si su nombre empieza con letra MAYÚSCULA.
type Usuario struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ==========================================
// 2. IMPLEMENTACIÓN POO (Estructura con métodos)
// ==========================================

// Nota: DB: Es una estructura contenedora que guarda el puntero a la conexión de la base de datos (*sql.DB). Esto nos permite adjuntarle métodos a la base de datos
type DB struct {
	db *sql.DB // Encapsulamos la conexión
}

// ==========================================
// 3. EJECUCIÓN DEL PROGRAMA (Main)
// ==========================================
func main() {
	// Reemplaza con tus credenciales reales
	// Cadena de conexión: usuario:contraseña@tcp(127.0.0.1:3306)/nombre_bdd
	var dsn = "root:@tcp(127.0.0.1:3306)/empresa"

	var db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error de configuración de BDD: %v", err)
	}
	defer db.Close()

	var app = &DB{db: db}

	// Servir archivos estáticos (HTML y CSS) desde la carpeta "public"
	var fs = http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Rutas de la API REST
	http.HandleFunc("/api/usuarios", app.handleUsuarios)
	http.HandleFunc("/api/usuarios/", app.handleUsuarioPorID)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Controlador para GET (todos) y POST (crear)
func (app *DB) handleUsuarios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		var rows, err = app.db.Query("SELECT id, name, email FROM usuarios")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var usuarios []Usuario
		for rows.Next() {
			var u Usuario
			rows.Scan(&u.Id, &u.Name, &u.Email)
			usuarios = append(usuarios, u)
		}
		json.NewEncoder(w).Encode(usuarios)

	case "POST":
		var u Usuario
		var errD = json.NewDecoder(r.Body).Decode(&u)
		if errD != nil {
			http.Error(w, errD.Error(), http.StatusBadRequest)
			return
		}
		var res, err = app.db.Exec("INSERT INTO usuarios (name, email) VALUES (?, ?)", u.Name, u.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var id, _ = res.LastInsertId()
		u.Id = int(id)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

// Controlador para PUT (actualizar) y DELETE (eliminar)
func (app *DB) handleUsuarioPorID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var idStr = r.URL.Path[len("/api/usuarios/"):]
	var id, errA = strconv.Atoi(idStr)
	if errA != nil {
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
		_, err = app.db.Exec("UPDATE usuarios SET name = ?, email = ? WHERE id = ?", u.Name, u.Email, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u.Id = id
		json.NewEncoder(w).Encode(u)

	case "DELETE":
		var _, err = app.db.Exec("DELETE FROM usuarios WHERE id = ?", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
