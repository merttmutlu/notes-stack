# Architecture Notes

## Core Entity

The system is centered around a single entity: `Note`.

`Note` fields:

- `id`
- `title`
- `content`
- `created_at`
- `updated_at`

This intentionally minimal model is enough to support:

- creating notes
- listing notes line by line
- viewing a single note
- editing notes
- deleting notes

The API, database schema, and frontend state should all map to this same model.

## API Technology Choices

The backend API is intentionally kept minimal to reduce maintenance overhead.

Standard library usage:

- `net/http` for the HTTP server
- `encoding/json` for JSON encoding and decoding
- `log/slog` for structured logging
- `context` for request-scoped operations

External dependency:

- `github.com/jackc/pgx/v5` at `v5.8.0`

Reasoning:

- Go does not include a PostgreSQL driver in the standard library.
- `pgx` is the only required external dependency for database connectivity.
- The project avoids web frameworks, ORMs, validation libraries, and extra config libraries in the first version.

## API Implementation Plan

The first API implementation should focus on the smallest end-to-end path between HTTP requests and the `notes` table.

Initial package layout:

- `apps/api/cmd/api` for the application entrypoint
- `apps/api/internal/config` for environment-based configuration
- `apps/api/internal/http` for route registration and handlers
- `apps/api/internal/store` for PostgreSQL queries
- `apps/api/internal/note` for note-specific types

Implementation order:

1. Create the Go module and API entrypoint.
2. Load configuration from environment variables.
3. Connect the API to PostgreSQL using `pgx`.
4. Implement `GET /healthz`.
5. Implement `GET /readyz` with a database connectivity check.
6. Implement `GET /api/v1/notes` to return existing rows.
7. Implement `POST /api/v1/notes` to insert a new note.
8. Add `GET /api/v1/notes/{id}`, `PUT`, and `DELETE` after the basic flow works.

First API milestone:

- the server starts
- health endpoints respond correctly
- the API can query the `notes` table
- the API can create and list notes
