FROM golang:alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd/api
CMD ["app"]
EXPOSE 4000
