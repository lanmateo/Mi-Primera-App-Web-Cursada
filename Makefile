.PHONY: test

test:
	@echo "--- 1) TAREAS PREVIAS ---"
	@echo "Ejecutando sqlc..."
	sqlc generate
	@echo "Compilando proyecto..."
	go build ./...
	@echo "Borrando contenedores y volúmenes viejos..."
	docker-compose down -v
	@echo "Levantando contenedor de base de datos..."
	docker-compose up -d db-test
	@echo "Esperando a que la base de datos esté corriendo..."
	@until docker inspect -f '{{.State.Health.Status}}' $$(docker-compose ps -q db-test) | grep -q "healthy"; do \
		sleep 1; \
	done
	@echo "Base de datos lista."

	@echo "\n--- 2) EJECUCIÓN DE TESTS ---"
	go test ./db/... -v

	@echo "\n--- 3) TAREAS POSTERIORES ---"
	@echo "Borrando contenedores y volúmenes..."
	docker-compose down -v
