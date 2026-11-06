#!/usr/bin/env bash
# Ejecuta la aplicación completa:
#   1) Tareas previas: .env, sqlc, templ, compilación, limpieza y arranque de la BD en Docker
#   2) Ejecución de la aplicación (http://localhost:8080)
#   3) Tareas posteriores: detener la app y borrar contenedores y volúmenes (al salir con Ctrl+C)
set -euo pipefail

# Trabajar siempre desde la carpeta del proyecto
cd "$(dirname "$0")"

LIMPIEZA_HECHA=0
cleanup() {
  if [ "$LIMPIEZA_HECHA" -eq 1 ]; then return; fi
  LIMPIEZA_HECHA=1
  echo ""
  echo "--- 3) TAREAS POSTERIORES ---"
  if [ -n "${SERVER_PID:-}" ]; then
    echo "Deteniendo la aplicación (PID $SERVER_PID)..."
    kill "$SERVER_PID" 2>/dev/null || true
  fi
  echo "Borrando contenedores y volúmenes..."
  docker compose down -v
}
trap cleanup EXIT
trap 'exit 0' INT TERM

echo "--- 1) TAREAS PREVIAS ---"
if [ ! -f .env ]; then
  echo "Creando .env a partir de env.example..."
  cp env.example .env
fi
# Exportar las variables de .env para que las use la aplicación
set -a
source .env
set +a

echo "Ejecutando sqlc..."
if command -v sqlc >/dev/null 2>&1; then
  sqlc generate
else
  echo "(sqlc no está instalado: se usa la imagen de Docker sqlc/sqlc)"
  docker run --rm -u "$(id -u):$(id -g)" -v "$PWD":/src -w /src sqlc/sqlc generate
fi

echo "Ejecutando templ generate..."
# Se usa la misma versión de templ que figura en go.mod
TEMPL_VERSION=$(go list -m -f '{{.Version}}' github.com/a-h/templ)
go run github.com/a-h/templ/cmd/templ@"$TEMPL_VERSION" generate

echo "Compilando la aplicación..."
mkdir -p tmp
go build -o tmp/app .

echo "Borrando contenedores y volúmenes viejos..."
docker compose down -v

echo "Levantando contenedor de base de datos..."
docker compose up -d db-test

echo "Esperando a que la base de datos esté corriendo..."
until docker inspect -f '{{.State.Health.Status}}' "$(docker compose ps -q db-test)" 2>/dev/null | grep -q "healthy"; do
  sleep 1
done
echo "Base de datos lista."

echo ""
echo "--- 2) EJECUCIÓN DE LA APLICACIÓN ---"
echo "Abrí http://localhost:${PORT:-8080} en el navegador."
echo "Presioná Ctrl+C para detener la aplicación y limpiar los contenedores."
./tmp/app &
SERVER_PID=$!
wait "$SERVER_PID"
