# TP3: Exponiendo la Lógica como API - Catálogo de Películas

## Descripción del Proyecto
Este proyecto implementa una API REST construida en Go para gestionar un catálogo de películas, conectando la lógica de la aplicación con una capa de persistencia en PostgreSQL gestionada mediante `sqlc`.

## Endpoints de la API

| Método | Endpoint | Descripción | Códigos HTTP esperados |
| :--- | :--- | :--- | :--- |
| `POST` | `/peliculas` | Crea una nueva película a partir de un cuerpo JSON. Valida que el título no esté vacío. | `201 Created`, `400 Bad Request` |
| `GET` | `/peliculas` | Lista todas las películas almacenadas en formato JSON. | `200 OK` |
| `GET` | `/peliculas/{id}` | Obtiene los detalles de una película según su identificador. | `200 OK`, `404 Not Found` |
| `PUT` | `/peliculas/{id}` | Actualiza los datos de una película existente mediante JSON. | `200 OK`, `400 Bad Request`, `404 Not Found` |
| `DELETE` | `/peliculas/{id}` | Elimina una película por su identificador. | `204 No Content`, `404 Not Found` |

### Formato de Entrada JSON (`POST` y `PUT`)
```json
{
  "titulo": "Inception",
  "descripcion": "Ladrones de sueños",
  "duracion": 148,
  "actores": "Leonardo DiCaprio",
  "genero": "Ciencia Ficción",
  "anio_lanzamiento": 2010,
  "puntaje": 8.8
}
```
## Cómo ejecutar las pruebas de la API
Para probar la API de forma automatizada, clonar el repositorio, ubicarse en la rama tp3 y ejecutar:
make test-api

(Alternativamente, puede ejecutarse directamente el script ./test-api.sh).

Este comando realiza de manera autónoma el ciclo completo:

1. Tareas Previas: Genera el archivo .env a partir de env.example (si no existe), genera el código Go con sqlc, compila el binario del servidor, limpia contenedores previos, levanta PostgreSQL en Docker con docker compose, espera a que esté saludable e inicia el servidor web de la API.

2. Prueba de Endpoints: Ejecuta la batería de pruebas sobre cada endpoint validando los códigos de respuesta (201, 200, 204, 400, 404).

3. Tareas Posteriores: Detiene el servidor web de la aplicación y destruye los contenedores y volúmenes creados con docker compose down -v.

Las peticiones de prueba individuales también se encuentran documentadas en los archivos requests.sh (cURL) y requests.hurl (HURL).
