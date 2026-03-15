# Note Data Model

The core entity in `notes-stack` is `Note`.

## Purpose

A note is the single business object stored by the application. The first version of the project supports creating notes, listing them, viewing a single note, editing them, and deleting them.

## Fields

| Field | Type | Required | Set By | Description |
| --- | --- | --- | --- | --- |
| `id` | UUID | Yes | Server | Unique identifier for a note. Used to fetch, update, and delete a specific record. |
| `title` | string | Yes | Client | Short text shown in the notes list. |
| `content` | string | Yes | Client | Main body of the note. |
| `created_at` | timestamp | Yes | Server | Time when the note was first created. |
| `updated_at` | timestamp | Yes | Server | Time when the note was last modified. |

## Representation

Example JSON representation:

```json
{
  "id": "3d4ca431-7db8-4f5e-868f-5d8b5fd5c737",
  "title": "Check Kubernetes logs",
  "content": "Use kubectl logs and describe to inspect the failing pod.",
  "created_at": "2026-03-15T10:00:00Z",
  "updated_at": "2026-03-15T10:00:00Z"
}
```

## Create Input

When creating a note, the client only sends the user-controlled fields:

```json
{
  "title": "Check Kubernetes logs",
  "content": "Use kubectl logs and describe to inspect the failing pod."
}
```

The server is responsible for generating `id`, `created_at`, and `updated_at`.

## Initial Rules

- `title` is required.
- `content` is required.
- Notes are listed line by line in the UI.
- The server owns identifiers and timestamps.
- The first version keeps a single entity and does not include users, tags, or categories.

## API Plan

The API will expose health endpoints and note CRUD endpoints.

Planned endpoints:

- `GET /healthz`
- `GET /readyz`
- `GET /api/v1/notes`
- `GET /api/v1/notes/{id}`
- `POST /api/v1/notes`
- `PUT /api/v1/notes/{id}`
- `DELETE /api/v1/notes/{id}`

Initial implementation priority:

1. `GET /healthz`
2. `GET /readyz`
3. `GET /api/v1/notes`
4. `POST /api/v1/notes`

The first implementation should use Go standard library packages for HTTP and JSON, and `github.com/jackc/pgx/v5` at `v5.8.0` for PostgreSQL access.

## Endpoint Contract

### `GET /healthz`

Purpose:

- confirms that the API process is running

Response:

- `200 OK`

```json
{
  "status": "ok"
}
```

### `GET /readyz`

Purpose:

- confirms that the API is ready to serve traffic
- should include a PostgreSQL connectivity check

Success response:

- `200 OK`

```json
{
  "status": "ready"
}
```

Failure response:

- `503 Service Unavailable`

```json
{
  "status": "not_ready"
}
```

### `GET /api/v1/notes`

Purpose:

- returns all notes from the database

Behavior:

- list notes in descending `created_at` order

Success response:

- `200 OK`

```json
{
  "data": [
    {
      "id": "3d4ca431-7db8-4f5e-868f-5d8b5fd5c737",
      "title": "Check Kubernetes logs",
      "content": "Use kubectl logs and describe to inspect the failing pod.",
      "created_at": "2026-03-15T10:00:00Z",
      "updated_at": "2026-03-15T10:00:00Z"
    }
  ]
}
```

Empty list response:

- `200 OK`

```json
{
  "data": []
}
```

### `POST /api/v1/notes`

Purpose:

- creates a new note

Request body:

```json
{
  "title": "Check Kubernetes logs",
  "content": "Use kubectl logs and describe to inspect the failing pod."
}
```

Validation rules:

- `title` is required
- `content` is required

Success response:

- `201 Created`

```json
{
  "data": {
    "id": "3d4ca431-7db8-4f5e-868f-5d8b5fd5c737",
    "title": "Check Kubernetes logs",
    "content": "Use kubectl logs and describe to inspect the failing pod.",
    "created_at": "2026-03-15T10:00:00Z",
    "updated_at": "2026-03-15T10:00:00Z"
  }
}
```

Validation error response:

- `400 Bad Request`

```json
{
  "error": "title and content are required"
}
```

Server error response:

- `500 Internal Server Error`

```json
{
  "error": "internal server error"
}
```
