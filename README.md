# Olx API

A minimal Go API for managing OLX-style listings stored in PostgreSQL.

## Overview

This project provides a lightweight HTTP API for reading and removing marketplace listings. It is built around a Go HTTP server, a PostgreSQL database, and a simple migration-based schema setup. The current implementation supports health checks and listing operations, making it suitable as a starting point for a marketplace backend or a small internal API.

The service is intended to solve the problem of exposing a simple, database-backed API for listing data without introducing unnecessary framework complexity or infrastructure overhead.

## Features

- Health check endpoint for service availability
- Listing retrieval endpoint for recent marketplace records
- Listing deletion by identifier
- PostgreSQL-backed persistence
- Database migration support for schema setup and rollback
- Environment-based configuration via `.env` files

## Tech Stack

- Go
- PostgreSQL
- `pgx` for database connectivity
- `golang-migrate` for schema migrations
- `godotenv` for environment variables
- Standard library HTTP server and JSON handling

## Architectur

The application is organized as a small layered service:

- `cmd/api` starts the HTTP server and wires together the config, database, and handlers
- `internals/config` loads required runtime configuration from environment variables
- `internals/db` opens and validates the PostgreSQL connection
- `internals/handlers` exposes the HTTP routes and business logic
- `migrations` contains SQL migration files for creating and altering the database schema

```mermaid
flowchart LR
    Client[Client / Consumer] -->|HTTP| API[Go API Server]
    API --> Config[config.MustLoad]
    API --> DB[PostgreSQL Database]
    API --> Handlers[Listing & Health Handlers]
    Migrations[Migration Files] --> DB
```

The request flow is straightforward: the server loads configuration, opens a PostgreSQL connection, registers route handlers, and executes database queries or deletes based on the incoming request.

## Project Structure

```text
.
├── bin/
├── cmd/
│   ├── api/
│   │   └── main.go
│   └── migrate/
│       └── main.go
├── internals/
│   ├── config/
│   │   └── config.go
│   ├── db/
│   │   └── db.go
│   └── handlers/
│       ├── healthz.go
│       └── listings.go
├── migrations/
│   ├── 000001_create_listings_table.down.sql
│   └── 000001_create_listings_table.up.sql
├── README.md
├── go.mod
└── .env.example
```

## Current Status

This project is currently in active development and functions as a small MVP for listing management.

## Prerequisites

Before running the project locally, ensure the following are available:

- Go 1.26.6 or newer
- PostgreSQL database instance
- Access to a local or remote PostgreSQL server
- Git

## Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/aasif-10/Olx-API.git
   cd Olx-API
   ```

2. Create a `.env` file in the project root and configure the required environment variables.

3. Apply the database migration:

   ```bash
   go run ./cmd/migrate up
   ```

4. Start the API server:

   ```bash
   go run ./cmd/api
   ```

## Environment Variables

The application expects the following environment variables:

```env
PORT=8080
ENV=development
DATABASE_URL=postgres://username:password@localhost:5432/olx_api
```

| Variable       | Required | Description                                                     |
| -------------- | -------- | --------------------------------------------------------------- |
| `PORT`         | Yes      | Port used by the HTTP server.                                   |
| `ENV`          | Yes      | Runtime environment name such as `development` or `production`. |
| `DATABASE_URL` | Yes      | PostgreSQL connection string used by the application.           |

## Running the Project

### Development

```bash
go run ./cmd/api
```

### Database migration

Apply schema changes:

```bash
go run ./cmd/migrate up
```

Revert the last migration:

```bash
go run ./cmd/migrate down
```

### Production

This project does not currently define a production deployment workflow in the repository. Use the configured environment variables and host process that matches your deployment platform.

## API Documentation

### GET /healthz

Returns the service status.

- Method: `GET`
- Endpoint: `/healthz`
- Purpose: Check whether the API is running
- Parameters: None

Example request:

```bash
curl http://localhost:8080/healthz
```

Example response:

```json
{ "status": "all ok" }
```

Possible errors:

- `500` is not expected for the current implementation; the handler currently returns a successful response with no error path.

### GET /listings

Returns recently created listings ordered by `created_at` descending, limited to 100 records.

- Method: `GET`
- Endpoint: `/listings`
- Purpose: Retrieve marketplace listing records
- Parameters: None

Example request:

```bash
curl http://localhost:8080/listings
```

Example response:

```json
[
  {
    "id": "6f6d8c71-3d7e-4c14-b9da-0d9d0f53d1b8",
    "title": "Apartment for rent",
    "description": "2-bedroom apartment in the city center",
    "price": "350000",
    "city": "Lagos",
    "created_at": "2026-08-30T10:00:00Z"
  }
]
```

Possible errors:

- `500 Internal Server Error` if the database query fails or rows cannot be scanned

### DELETE /listings/{id}

Deletes a listing by its identifier.

- Method: `DELETE`
- Endpoint: `/listings/{id}`
- Purpose: Remove a listing record
- Parameters:
  - `id` (path parameter, UUID)

Example request:

```bash
curl -X DELETE http://localhost:8080/listings/6f6d8c71-3d7e-4c14-b9da-0d9d0f53d1b8
```

Example response:

```http
HTTP/1.1 204 No Content
```

Possible errors:

- `500 Internal Server Error` if the delete query fails

## Database

The project uses PostgreSQL as its primary data store.

Current schema:

```sql
CREATE TABLE listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    price BIGINT NOT NULL,
    city TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

Migrations are stored under the `migrations/` directory and are managed with `golang-migrate`.

## Testing

There is no automated test suite in the repository at this time.

The current validation flow is manual and consists of:

- starting the API locally,
- confirming the database connection,
- calling the health and listing endpoints,
- verifying that delete operations work as expected.

## Deployment

Deployment instructions are not yet defined in this repository.

Use the following as the current deployment model:

- configure `PORT`, `ENV`, and `DATABASE_URL` in the target environment,
- run the API process with the appropriate host/runtime environment,
- ensure PostgreSQL is reachable from the deployment environment,
- apply migrations before serving traffic.

## Design Decisions / Technical Details

- The service uses the Go standard library HTTP server instead of a framework to keep the implementation small and explicit.
- Database connections are managed through `database/sql` with `pgx` as the driver and basic pooling settings.
- Configuration is loaded from environment variables at startup, avoiding hard-coded runtime values.
- The API returns JSON responses for both health and listing operations.
- Listing queries are limited to the most recent 100 rows to prevent unbounded result payloads.

## Future Improvements

The repository does not yet define a specific roadmap. Potential future work would be driven by requirements added to the project later, but no formal improvement plan is currently included.

## Contributing

Contributions are welcome. Please keep changes focused and consistent with the current project structure.

Suggested workflow:

1. Create a feature branch from the main branch.
2. Make the smallest clear change that addresses the problem.
3. Validate the behavior with the relevant local checks.
4. Open a pull request describing the change and any usage impact.

## License

This repository does not currently include a license file. Before publishing or distributing the project, add an appropriate open-source license such as MIT or Apache-2.0.

## Repository

- GitHub: https://github.com/aasif-10/Olx-API
