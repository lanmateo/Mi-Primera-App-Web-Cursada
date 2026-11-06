# TP5: Renderizado del Lado del Servidor con Templ - Catálogo de Películas

## Descripción del Proyecto
Aplicación web en Go para gestionar un catálogo de películas. A partir de este TP se abandona el enfoque **API REST + cliente JavaScript**: el HTML se genera completo en el servidor (SSR) usando componentes **templ**, y los datos se obtienen de **PostgreSQL** mediante el código generado por **sqlc**.

La aplicación no utiliza archivos JavaScript: el alta y la baja de películas se hacen con formularios HTML clásicos y el patrón **Post/Redirect/Get (PRG)**.

## Cómo ejecutar la aplicación

### Requisitos
- Go 1.23 o superior (si la versión instalada es anterior, Go descarga automáticamente la necesaria)
- Docker con Docker Compose
- `sqlc` (opcional: si no está instalado, el script usa la imagen de Docker `sqlc/sqlc`)
- `templ` **no** necesita instalarse: el script lo ejecuta con `go run`, usando la versión fijada en `go.mod`

### Ejecución
Clonar el repositorio, ubicarse en la rama `tp5` y ejecutar:

```bash
make run-app
```

*(Alternativamente: `./run-app.sh`)*

Luego abrir **http://localhost:8080** en el navegador. Para terminar, presionar **Ctrl+C**.

El script realiza automáticamente:
1. **Tareas previas:** crea el `.env` a partir de `env.example` (si no existe), ejecuta `sqlc generate`, ejecuta `templ generate`, compila la aplicación, borra contenedores y volúmenes previos, levanta PostgreSQL con `docker compose` y espera a que la base esté saludable.
2. **Ejecución de la aplicación:** inicia el servidor web en el puerto `8080`.
3. **Tareas posteriores:** al presionar Ctrl+C (o si la aplicación termina), detiene el servidor y borra contenedores y volúmenes con `docker compose down -v`.

### Tests de persistencia
```bash
make test
```
Levanta la base en Docker, ejecuta los tests de `db/sqlc/db_test.go` y luego borra contenedores y volúmenes.

## Rutas de la aplicación

| Método | Ruta | Descripción |
| :--- | :--- | :--- |
| `GET` | `/` | Página principal: formulario de alta y tabla con todas las películas. |
| `POST` | `/peliculas` | Recibe el formulario, crea la película y redirige a `/` (PRG). |
| `POST` | `/peliculas/{id}/eliminar` | Elimina la película y redirige a `/` (PRG). |

Si algún dato del formulario es inválido (título vacío, duración, año o puntaje no numéricos), se responde `400 Bad Request`. El formulario además valida los campos en el navegador con atributos HTML (`required`, `min`, `max`, `step`).

---

## Documentación de Arquitectura

```
Navegador ──GET /──────────────▶ handlers.Index  ──ListPeliculas──▶ sqlc ──▶ PostgreSQL
          ◀──HTML completo──────  views.IndexPage (templ)

Navegador ──POST /peliculas────▶ handlers.Create ──CreatePelicula─▶ sqlc ──▶ PostgreSQL
          ◀──303 See Other: / ──  (el navegador vuelve a pedir GET /)
```

### Capa de Presentación (`views/`, componentes templ)
- `layout.templ` → `Layout`: estructura base de la página (`<!DOCTYPE html>`, `<head>`, `<body>`), con el enlace a **Pico.css**. Recibe el contenido de cada página como hijo (`{ children... }`).
- `entity_form.templ` → `PeliculaForm`: formulario de alta que hace `POST /peliculas`.
- `entity_list.templ` → `PeliculaList`: recibe `[]db.Pelicula` y renderiza la tabla (o un mensaje si no hay películas). Cada fila tiene un formulario `POST /peliculas/{id}/eliminar`.
- `index.templ` → `IndexPage`: página principal; usa `Layout` y renderiza el título, el formulario y la lista.

Los archivos `*_templ.go` se generan con `templ generate` y no se versionan (están en `.gitignore`).

### Capa de Handlers (`handlers/peliculas.go`)
- `Index` (`GET /`): obtiene las películas con `ListPeliculas` y renderiza `views.IndexPage` en el `http.ResponseWriter`.
- `Create` (`POST /peliculas`): parsea el formulario con `r.ParseForm()`, convierte los campos numéricos, guarda con `CreatePelicula` y redirige con `http.Redirect(w, r, "/", http.StatusSeeOther)`.
- `Delete` (`POST /peliculas/{id}/eliminar`): borra con `DeletePelicula` y redirige a `/`.

### `main.go`
Conecta con PostgreSQL usando las variables de entorno, crea el repositorio de sqlc (`db.New`), registra las rutas en un `http.ServeMux` y levanta el servidor en el puerto `8080`.

### Capa de Persistencia (`sqlc` + PostgreSQL)
Se mantiene la capa del TP2/TP3: **sqlc** genera funciones Go tipadas a partir de SQL puro.

#### Tabla `peliculas` (`db/schema/schema.sql`)
- `id` (`SERIAL`, Primary Key)
- `titulo` (`VARCHAR(255)`)
- `descripcion` (`TEXT`)
- `duracion` (`INT`)
- `actores` (`VARCHAR(255)`)
- `genero` (`VARCHAR(100)`)
- `anio_lanzamiento` (`INT`)
- `puntaje` (`FLOAT`)

#### Consultas (`db/queries/queries.sql`)
- `CreatePelicula`, `GetPelicula`, `ListPeliculas`, `UpdatePelicula`, `DeletePelicula`.

### Cambios respecto del TP3
- Se eliminaron los handlers que devolvían JSON y los scripts de prueba de la API (`test-api.sh`, `requests.sh`, `requests.hurl`).
- Se agregaron los componentes templ en `views/` y el script `run-app.sh` (`make run-app`).
- El healthcheck de PostgreSQL usa `pg_isready -h localhost`, para que la base figure como lista recién cuando acepta conexiones TCP.

## Estructura del proyecto
```
.
├── main.go
├── handlers/
│   └── peliculas.go
├── views/
│   ├── layout.templ
│   ├── entity_form.templ
│   ├── entity_list.templ
│   └── index.templ
├── db/
│   ├── schema/schema.sql
│   ├── queries/queries.sql
│   └── sqlc/db_test.go        (el resto de db/sqlc lo genera sqlc)
├── sqlc.yaml
├── docker-compose.yml
├── env.example
├── run-app.sh
├── Makefile
├── go.mod / go.sum
└── README.md
```
