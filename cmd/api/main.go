package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"crud/internal/handlers"
	"crud/internal/repository"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/empresa"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error de conexión: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	// 1. Instanciar Repositorios
	usuarioRepo := repository.NewMySQLUsuarioRepository(db)
	// productoRepo := repository.NewMySQLProductoRepository(db)

	// 2. Instanciar Handlers
	usuarioHandler := handlers.NewUsuarioHandler(usuarioRepo)

	// 3. Servir Frontend Estático
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// 4. Endpoints de la API
	http.HandleFunc("/api/usuarios", usuarioHandler.HandleUsuarios)
	http.HandleFunc("/api/usuarios/", usuarioHandler.HandleUsuarioPorID)

	fmt.Println("Servidor escuchando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
