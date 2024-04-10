build:
	@docker build -t dbasik:latest .

run-container:
	@docker run -it --rm -p 4000:4000 dbasik:latest

run:
	@docker compose up -d
