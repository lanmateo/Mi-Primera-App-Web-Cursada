package db

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq" // Driver de Postgres
)

func TestPersistenciaPeliculas(t *testing.T) {
	// 1. Conectarnos a la base de datos de prueba que levantó Docker
	connStr := "postgres://user:password@localhost:5432/peliculas_test?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("No se pudo conectar a la BD: %v", err)
	}
	defer db.Close()

	// 2. Inicializar el objeto Queries generado por sqlc
	queries := New(db)
	ctx := context.Background()

	// 3. Prueba de Create (Insertar una película)
	nuevaPelicula := CreatePeliculaParams{
		Titulo:          "Inception",
		Descripcion:     "Un ladrón que roba secretos corporativos a través del uso de la tecnología de compartir sueños.",
		Duracion:        148,
		Actores:         "Leonardo DiCaprio, Joseph Gordon-Levitt",
		Genero:          "Ciencia Ficción",
		AnioLanzamiento: 2010,
		Puntaje:         8.8,
	}

	peliculaInsertada, err := queries.CreatePelicula(ctx, nuevaPelicula)
	if err != nil {
		t.Fatalf("Fallo al insertar la pelicula: %v", err)
	}

	if peliculaInsertada.Titulo != nuevaPelicula.Titulo {
		t.Errorf("Se esperaba el titulo %s, se obtuvo %s", nuevaPelicula.Titulo, peliculaInsertada.Titulo)
	}

	// 4. Prueba de Get (Leer la película insertada)
	peliculaLeida, err := queries.GetPelicula(ctx, peliculaInsertada.ID)
	if err != nil {
		t.Fatalf("Fallo al obtener la pelicula: %v", err)
	}

	if peliculaLeida.ID != peliculaInsertada.ID {
		t.Errorf("Se esperaba el ID %d, se obtuvo %d", peliculaInsertada.ID, peliculaLeida.ID)
	}

    // 5. Prueba de Update (Actualizar la película)
	updateParams := UpdatePeliculaParams{
		ID:              peliculaInsertada.ID,
		Titulo:          "Inception editada",
		Descripcion:     peliculaInsertada.Descripcion,
		Duracion:        150,
		Actores:         peliculaInsertada.Actores,
		Genero:          peliculaInsertada.Genero,
		AnioLanzamiento: peliculaInsertada.AnioLanzamiento,
		Puntaje:         9.0,
	}    
    err = queries.UpdatePelicula(ctx, updateParams)
	if err != nil {
		t.Fatalf("Fallo al actualizar la pelicula: %v", err)
	}

    // 6. Prueba de List (Listar todas las películas)
	peliculas, err := queries.ListPeliculas(ctx)
	if err != nil {
		t.Fatalf("Fallo al listar las peliculas: %v", err)
	}
	if len(peliculas) == 0 {
		t.Errorf("Se esperaba al menos 1 pelicula en la lista, se obtuvieron 0")
	}

    // 7. Prueba de Delete (Borrar la película)
	err = queries.DeletePelicula(ctx, peliculaInsertada.ID)
	if err != nil {
		t.Fatalf("Fallo al borrar la pelicula: %v", err)
	}

	t.Logf("Test completado con éxito. Pelicula probada: %s", peliculaLeida.Titulo)
}
