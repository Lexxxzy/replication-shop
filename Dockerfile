FROM golang:1.21-alpine as builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" go build -o main cmd/api/main.go

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/main /app/.env /app/*.toml ./
ENTRYPOINT ["/app/main"]