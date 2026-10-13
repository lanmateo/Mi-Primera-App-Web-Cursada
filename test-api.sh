#!/usr/bin/env bash
set -e

# Si ocurre cualquier error o interrupción, limpiar siempre los recursos
cleanup() {
  echo -e "\n--- 3) TAREAS POSTERIORES ---"
  if [ -n "$SERVER_PID" ]; then
    echo "Deteniendo servidor de la API (PID $SERVER_PID)..."
    kill "$SERVER_PID" 2>/dev/null || true
  fi
  echo "Borrando contenedores y volúmenes de Docker..."
  docker compose down -v
}
trap cleanup EXIT INT TERM

echo "--- 1) TAREAS PREVIAS ---"
if [ ! -f .env ]; then
  echo "Creando .env a partir de env.example..."
  cp env.example .env
fi

echo "Ejecutando sqlc..."
sqlc generate

echo "Compilando proyecto..."
mkdir -p tmp
go build -o tmp/api-server main.go

echo "Borrando contenedores y volúmenes viejos..."
docker compose down -v

echo "Levantando contenedor de base de datos..."
docker compose up -d db-test

echo "Esperando a que la base de datos esté saludable..."
until docker inspect -f '{{.State.Health.Status}}' $(docker compose ps -q db-test) 2>/dev/null | grep -q "healthy"; do
  sleep 1
done
echo "Base de datos lista."

echo "Iniciando servidor API en segundo plano..."
./tmp/api-server &
SERVER_PID=$!

echo "Esperando a que el servidor web responda en puerto 8080..."
until curl -s http://localhost:8080/peliculas >/dev/null 2>&1; do
  sleep 0.5
done
echo "Servidor web listo."

echo -e "\n--- 2) PRUEBA DE ENDPOINTS ---"
./requests.sh
