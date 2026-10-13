# 1. Crear película
POST http://localhost:8080/peliculas
Content-Type: application/json
{
  "titulo": "Inception",
  "descripcion": "Ladrones de sueños",
  "duracion": 148,
  "actores": "Leonardo DiCaprio",
  "genero": "Ciencia Ficción",
  "anio_lanzamiento": 2010,
  "puntaje": 8.8
}
HTTP 201
[Asserts]
jsonpath "$.id" isInteger
jsonpath "$.titulo" == "Inception"

# 2. Listar todas las películas
GET http://localhost:8080/peliculas
HTTP 200
[Asserts]
jsonpath "$.length()" >= 1

# 3. Obtener película por ID
GET http://localhost:8080/peliculas/1
HTTP 200
[Asserts]
jsonpath "$.titulo" == "Inception"

# 4. Actualizar película
PUT http://localhost:8080/peliculas/1
Content-Type: application/json
{
  "titulo": "Inception (Edicion Especial)",
  "descripcion": "Version extendida",
  "duracion": 155,
  "actores": "Leonardo DiCaprio, Elliot Page",
  "genero": "Ciencia Ficción",
  "anio_lanzamiento": 2010,
  "puntaje": 9.0
}
HTTP 200
[Asserts]
jsonpath "$.titulo" == "Inception (Edicion Especial)"
jsonpath "$.puntaje" == 9.0

# 5. Eliminar película
DELETE http://localhost:8080/peliculas/1
HTTP 204

# 6. Verificar que ya no existe
GET http://localhost:8080/peliculas/1
HTTP 404
