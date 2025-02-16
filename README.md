# Provisioning Server

A Go REST API for managing server fleets with admin authentication.

## Development Setup

### Prerequisites
- Go 1.21+
- Docker (optional)
- Goose
    ```
    go install github.com/pressly/goose/v3/cmd/goose@latest
    ```


### With Docker
```bash
# Start services
docker-compose up --build

# Run migrations
make migrate-up

# Stop and cleanup
docker-compose down -v
```

### Without Docker
```bash
export JWT_SECRET=...
export DB_CONNECTION_STRING=postgres://...

# Migrate db
make migrate-up

# Start server
go run main.go
```

### Standards



## API Endpoints

| Method | Endpoint           | Description                | Auth Required |
|--------|--------------------|----------------------------|---------------|
| POST   | /admin/register    | Register admin             | No            |
| POST   | /admin/login       | Login admin                | No            |
| POST   | /servers/register  | Register server            | No            |
| POST   | /servers           | List servers               | Admin         |
| POST   | /servers/approve   | Approve server             | Admin         |
