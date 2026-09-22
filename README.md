# Book Management System

A small REST API for managing books, built with Go, Gin, GORM, and MySQL. The project is designed as a learning reference for building a database-backed Go API with explicit dependencies, automated migrations, tests, and Docker Compose.

## Features

- Create, list, retrieve, update, and delete books
- MySQL persistence through GORM v2
- JSON request and response bodies
- Validation for malformed JSON and invalid book IDs
- In-memory SQLite tests that do not require MySQL
- Docker Compose setup for the API and database

## Requirements

- Go 1.27.1 or newer
- Docker Engine with Docker Compose, if using the containerized setup
- MySQL 8 or newer, if running the Go process directly

## Run With Docker Compose

This is the easiest way to run the complete application. From the project directory:

```bash
docker compose up --build
```

The API is available at <http://localhost:8080>. Press `Ctrl+C` to stop both the API and MySQL containers.

To stop containers started in the background:

```bash
docker compose down
```

The database is stored in the `mysql_data` Docker volume. To remove the containers and all stored database data:

```bash
docker compose down -v
```

The credentials in `compose.yaml` are for local development only. Do not use them in production.

## Run Locally

Start a local MySQL server and create the database and user:

```sql
CREATE DATABASE simplerest;
CREATE USER 'kyle'@'localhost' IDENTIFIED BY 'changeMe12';
GRANT ALL PRIVILEGES ON simplerest.* TO 'kyle'@'localhost';
FLUSH PRIVILEGES;
```

Then set the connection string and start the API:

```bash
export DATABASE_DSN='kyle:changeMe12@tcp(127.0.0.1:3306)/simplerest?charset=utf8mb4&parseTime=True&loc=Local'
go run ./cmd/main
```

`DATABASE_DSN` is optional when using the default local development values in `pkg/config`. For shared machines, CI, and production, always provide it through the environment rather than editing source code.

## API Endpoints

| Method | Endpoint | Description |
| --- | --- | --- |
| `POST` | `/book/` | Create a book |
| `GET` | `/book/` | List all books |
| `GET` | `/book/:bookId` | Retrieve one book |
| `PUT` | `/book/:bookId` | Update supplied book fields |
| `DELETE` | `/book/:bookId` | Delete a book |

Example request:

```bash
curl -i -X POST http://localhost:8080/book/ \
  -H 'Content-Type: application/json' \
  -d '{"name":"Dune","author":"Frank Herbert","publication":"1965"}'
```

Example update:

```bash
curl -i -X PUT http://localhost:8080/book/1 \
  -H 'Content-Type: application/json' \
  -d '{"publication":"1965 edition"}'
```

A book has this JSON shape:

```json
{
  "id": 1,
  "name": "Dune",
  "author": "Frank Herbert",
  "publication": "1965"
}
```

Successful creation returns `201 Created`. Successful deletion returns `204 No Content`. Invalid input returns `400 Bad Request`, and a missing book returns `404 Not Found`.

## Project Structure

```text
cmd/main/                 Application entry point
pkg/config/               Database connection configuration
pkg/controllers/          HTTP handlers and request validation
pkg/models/               Book entity and GORM repository
pkg/routes/               HTTP route registration
compose.yaml              Local MySQL and API services
Dockerfile                Multi-stage API image build
```

The application wires dependencies in `main`: database connection, repository, handler, and routes. This keeps package initialization free of database side effects and allows the handlers to be tested with an in-memory database.

## Tests and Quality Checks

Run the test suite:

```bash
go test ./...
```

Run static analysis:

```bash
go vet ./...
```

The controller tests use SQLite in memory, so they do not need a running MySQL server.

## Troubleshooting

### `connect: connection refused`

MySQL is not accepting connections on the configured address. Start the full stack with `docker compose up --build`, or start your local MySQL service and verify `DATABASE_DSN`.

### `access denied for user`

The username or password in `DATABASE_DSN` does not match the MySQL user. Check the credentials and database permissions.

### Port already in use

Stop the process using port `3306` or `8080`, or change the host-side port mapping in `compose.yaml`.
