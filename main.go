package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Cadena de conexión: usuario:contraseña@tcp(127.0.0.1:3306)/nombre_bdd
	var dsn = "root:@tcp(127.0.0.1:3306)/empresa"

	var db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error de conexión: %v", err)
	}
	defer db.Close()

	// 1. Instanciamos el Repositorio (Capa de datos)
	var repoUser = NewMySQLUsuarioRepository(db)

	// 2. Inyectamos el Repositorio dentro del Handler (Capa HTTP)
	var handler = NewUsuarioHandler(repoUser)

	// 3. Servidor de archivos estáticos HTML/CSS
	var fs = http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// 4. Rutas asociadas a los métodos del Handler
	http.HandleFunc("/api/usuarios", handler.HandleUsuarios)
	http.HandleFunc("/api/usuarios/", handler.HandleUsuarioPorID)

	fmt.Println("Servidor corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
