build:
	@docker build -t dbasik:latest .

run-container:
	@docker run -it --rm -p 4000:4000 dbasik:latest

run:
	@docker compose up -d

migrate:
	docker exec -it dbasik-db-1 migrate -path=${DBASIK_DB_DSN} -database= up

migrate-down:
	docker exec -it dbasik-db-1 migrate -path=${DBASIK_DB_DSN} -database= down

