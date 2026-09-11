CREATE TABLE peliculas (
    id SERIAL PRIMARY KEY,
    titulo VARCHAR(255) NOT NULL,
    descripcion TEXT NOT NULL,
    duracion INT NOT NULL,
    actores VARCHAR(255) NOT NULL,
    genero VARCHAR(100) NOT NULL,
    anio_lanzamiento INT NOT NULL,
    puntaje FLOAT NOT NULL
);
