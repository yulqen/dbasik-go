build:
	@docker build -t dbasik:latest .

run:
	@docker run -it --rm -p 4000:4000 dbasik:latest
