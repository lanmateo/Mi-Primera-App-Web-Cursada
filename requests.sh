#!/usr/bin/env bash
set -e

BASE_URL="http://localhost:8080/peliculas"

echo "=========================================="
echo "    INICIANDO PRUEBAS DE LA API REST     "
echo "=========================================="

# 1. POST /peliculas (Crear película válida -> Espera 201)
echo -e "\n[1/8] POST /peliculas - Creando película válida..."
HTTP_CODE=$(curl -s -o /tmp/resp_create.json -w "%{http_code}" -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{"titulo":"Inception","descripcion":"Ladrones de sueños","duracion":148,"actores":"Leonardo DiCaprio","genero":"Ciencia Ficción","anio_lanzamiento":2010,"puntaje":8.8}')

if [ "$HTTP_CODE" -eq 201 ]; then
  echo "✔ Correcto (HTTP 201 Created)"
  cat /tmp/resp_create.json
  echo ""
else
  echo "✖ Falló: Esperado 201, recibido $HTTP_CODE"
  exit 1
fi

# 2. POST /peliculas (Validación título vacío -> Espera 400)
echo -e "\n[2/8] POST /peliculas - Validando título vacío..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -d '{"titulo":"","descripcion":"Sin titulo","duracion":90,"actores":"N/A","genero":"Drama","anio_lanzamiento":2020,"puntaje":5.0}')

if [ "$HTTP_CODE" -eq 400 ]; then
  echo "✔ Correcto (HTTP 400 Bad Request ante validación fallida)"
else
  echo "✖ Falló: Esperado 400, recibido $HTTP_CODE"
  exit 1
fi

# 3. GET /peliculas (Listar películas -> Espera 200)
echo -e "\n[3/8] GET /peliculas - Listando todas las películas..."
HTTP_CODE=$(curl -s -o /tmp/resp_list.json -w "%{http_code}" -X GET "$BASE_URL")

if [ "$HTTP_CODE" -eq 200 ]; then
  echo "✔ Correcto (HTTP 200 OK)"
  cat /tmp/resp_list.json
  echo ""
else
  echo "✖ Falló: Esperado 200, recibido $HTTP_CODE"
  exit 1
fi

# 4. GET /peliculas/1 (Obtener por ID -> Espera 200)
echo -e "\n[4/8] GET /peliculas/1 - Consultando película con ID 1..."
HTTP_CODE=$(curl -s -o /tmp/resp_get.json -w "%{http_code}" -X GET "$BASE_URL/1")

if [ "$HTTP_CODE" -eq 200 ]; then
  echo "✔ Correcto (HTTP 200 OK)"
  cat /tmp/resp_get.json
  echo ""
else
  echo "✖ Falló: Esperado 200, recibido $HTTP_CODE"
  exit 1
fi

# 5. GET /peliculas/999 (Recurso inexistente -> Espera 404)
echo -e "\n[5/8] GET /peliculas/999 - Consultando película inexistente..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/999")

if [ "$HTTP_CODE" -eq 404 ]; then
  echo "✔ Correcto (HTTP 404 Not Found)"
else
  echo "✖ Falló: Esperado 404, recibido $HTTP_CODE"
  exit 1
fi

# 6. PUT /peliculas/1 (Actualizar película existente -> Espera 200)
echo -e "\n[6/8] PUT /peliculas/1 - Actualizando datos de película con ID 1..."
HTTP_CODE=$(curl -s -o /tmp/resp_put.json -w "%{http_code}" -X PUT "$BASE_URL/1" \
  -H "Content-Type: application/json" \
  -d '{"titulo":"Inception (Edicion Especial)","descripcion":"Version extendida","duracion":155,"actores":"Leonardo DiCaprio, Elliot Page","genero":"Ciencia Ficción","anio_lanzamiento":2010,"puntaje":9.0}')

if [ "$HTTP_CODE" -eq 200 ]; then
  echo "✔ Correcto (HTTP 200 OK con datos actualizados)"
  cat /tmp/resp_put.json
  echo ""
else
  echo "✖ Falló: Esperado 200, recibido $HTTP_CODE"
  exit 1
fi

# 7. DELETE /peliculas/1 (Eliminar película -> Espera 204)
echo -e "\n[7/8] DELETE /peliculas/1 - Eliminando película con ID 1..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE "$BASE_URL/1")

if [ "$HTTP_CODE" -eq 204 ]; then
  echo "✔ Correcto (HTTP 204 No Content)"
else
  echo "✖ Falló: Esperado 204, recibido $HTTP_CODE"
  exit 1
fi

# 8. GET /peliculas/1 (Verificar eliminación -> Espera 404)
echo -e "\n[8/8] GET /peliculas/1 - Verificando eliminación..."
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X GET "$BASE_URL/1")

if [ "$HTTP_CODE" -eq 404 ]; then
  echo "✔ Correcto (HTTP 404 Not Found, registro confirmado como eliminado)"
else
  echo "✖ Falló: Esperado 404, recibido $HTTP_CODE"
  exit 1
fi

echo -e "\n=========================================="
echo "    TODOS LOS ENDPOINTS FUNCIONAN OK      "
echo "=========================================="
