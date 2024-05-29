FROM golang:alpine as build

WORKDIR /

COPY . .

RUN go build -o ./app ./cmd/dbasik-api 

FROM scratch
COPY --from=build /app /app
COPY --from=build /.env /.env
CMD ["/app"]
