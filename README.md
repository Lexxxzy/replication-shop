# Fault-Tolerant shop
Stack:
- Golang (Echo)
- pgpool-ii & PostgreSQL
- Redis
- Cassandra
- Nginx
- Docker

## Description
This is a simple shop application that is designed to be fault-tolerant. It is built using Golang and uses Echo as the web framework. The application uses PostgreSQL as the main database and Redis as the cache. It also uses Cassandra as the backup database. The application is deployed using Docker and is load balanced using Nginx. The application is also fault-tolerant as it uses pgpool-ii to handle failover and load balancing for PostgreSQL.

### Start the application
Start only postgres + app:
```bash
docker compose -f docker-compose.postgresql.yml up -d
docker compose -f docker-compose.app.yml up -d --build
```

Access the application at `http://localhost:80`

Start golang API without docker for development:
:warning: Make sure you have `go` installed on your machine.

1. Install `air` for hot reloading:
```bash
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s
```

2. Create dev.env file:
```bash
cp .env dev.env
```

3. Run the following commands:
```bash
go mod download
./bin/air
```


