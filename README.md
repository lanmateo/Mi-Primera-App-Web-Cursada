# TP2: Persistiendo el Dominio - Catálogo de Películas

## Descripción del Proyecto
Este proyecto consiste en el desarrollo del backend para un catálogo de películas. En esta etapa (TP2), se aisló e implementó la capa de persistencia de datos, preparando la estructura necesaria para almacenar la información de las películas de forma permanente.

## Documentación de Persistencia
Para la capa de acceso a datos, se optó por el enfoque de **"SQL Directo"** utilizando **PostgreSQL** como motor de base de datos. 

En lugar de utilizar un ORM tradicional, la interacción con la base de datos se gestiona mediante **sqlc**. Esta herramienta toma las consultas SQL puras definidas en el proyecto y compila funciones de Go fuertemente tipadas, evitando el uso de código repetitivo y previniendo inyecciones SQL de forma nativa.

### Estructura de la Base de Datos
Se definió la entidad principal en el archivo `db/schema/schema.sql` con la siguiente estructura:
- `id` (SERIAL, Primary Key)
- `titulo` (VARCHAR)
- `descripcion` (TEXT)
- `duracion` (INT)
- `actores` (VARCHAR)
- `genero` (VARCHAR)
- `anio_lanzamiento` (INT)
- `puntaje` (FLOAT)

### Operaciones CRUD
Las consultas básicas se encuentran definidas en `db/queries/queries.sql`, implementando:
- `CreatePelicula`: Inserta un nuevo registro.
- `GetPelicula`: Obtiene una película por su ID.
- `ListPeliculas`: Lista todas las películas ordenadas por título.
- `UpdatePelicula`: Actualiza los datos de un registro existente.
- `DeletePelicula`: Elimina una película por su ID.

El código Go autogenerado por `sqlc` a partir de estas consultas se aloja en el paquete `db/sqlc/`.

---

## Cómo ejecutar los tests

Para probar la capa de persistencia, se debe clonar el repositorio, posicionarse en la rama `tp2` y ejecutar:

```bash
make test
