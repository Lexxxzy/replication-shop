FROM golang:1.21-alpine as app

WORKDIR /app

COPY . .

RUN go build -o main cmd/api/main.go

ENTRYPOINT ["/app/main"]
