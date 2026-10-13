package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// Driver de PostgreSQL
	_ "github.com/lib/pq"

	// Paquetes internos
	db "catalogo-peliculas/db/sqlc"
	"catalogo-peliculas/handlers"
)

// getEnv obtiene una variable de entorno o devuelve un valor por defecto
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// connectDB establece y valida la conexión con PostgreSQL
func connectDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("POSTGRES_PORT", "5432")
		user := getEnv("POSTGRES_USER", "user")
		pass := getEnv("POSTGRES_PASSWORD", "password")
		name := getEnv("POSTGRES_DB", "peliculas_test")

		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, pass, name)
	}

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error al abrir configuración de BD: %w", err)
	}

	if err := dbConn.Ping(); err != nil {
		dbConn.Close()
		return nil, fmt.Errorf("error al conectar con la BD: %w", err)
	}

	return dbConn, nil
}

func main() {
	log.Println("Iniciando aplicación del Catálogo de Películas...")

	// 1. Conexión a la base de datos
	conn, err := connectDB()
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer conn.Close()
	log.Println("Conexión a PostgreSQL establecida.")

	// 2. Instancia del repositorio sqlc
	queries := db.New(conn)

	// 3. Registrar Handlers CRUD
	peliculaHandler := handlers.NewPeliculaHandler(queries)
	http.HandleFunc("/peliculas", peliculaHandler.HandlePeliculas)
	http.HandleFunc("/peliculas/", peliculaHandler.HandlePeliculaByID)

	port := getEnv("PORT", "8080")
	log.Printf("Servidor API escuchando en http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
