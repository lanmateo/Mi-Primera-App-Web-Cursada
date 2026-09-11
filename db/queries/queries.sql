-- name: CreatePelicula :one
INSERT INTO peliculas (titulo, descripcion, duracion, actores, genero, anio_lanzamiento, puntaje)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetPelicula :one
SELECT * FROM peliculas WHERE id = $1;

-- name: ListPeliculas :many
SELECT * FROM peliculas ORDER BY titulo;

-- name: UpdatePelicula :exec
UPDATE peliculas
SET titulo = $2, descripcion = $3, duracion = $4, actores = $5, genero = $6, anio_lanzamiento = $7, puntaje = $8
WHERE id = $1;

-- name: DeletePelicula :exec
DELETE FROM peliculas WHERE id = $1;
