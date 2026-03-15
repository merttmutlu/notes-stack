# notes-stack

`notes-stack` is a small three-tier practice project for Kubernetes, infrastructure, and CI/CD work.

The first application entity is `Note`, with this initial data model:

- `id: UUID`
- `title: string`
- `content: string`
- `created_at: timestamp`
- `updated_at: timestamp`

Detailed design notes live in [docs/api.md](/Users/mertmutlu/Documents/GitHub/notes-stack/docs/api.md).

The API design uses Go standard library packages where possible and a single external dependency for PostgreSQL access: `github.com/jackc/pgx/v5` at `v5.8.0`.

## Local Backend Flow

The local backend consists of:

- PostgreSQL running with Docker Compose
- a Go API running on port `8080`

### 1. Start PostgreSQL

```bash
task postgres:up
```

PostgreSQL data is stored under `deploy/compose/data`.

### 2. Apply Migrations

```bash
task postgres:migrate
```

This applies:

- [001_create_notes.sql](/Users/mertmutlu/Documents/GitHub/notes-stack/db/migrations/001_create_notes.sql)
- [002_enable_pgcrypto.sql](/Users/mertmutlu/Documents/GitHub/notes-stack/db/migrations/002_enable_pgcrypto.sql)

### 3. Run the API

```bash
task api:run
```

The API expects:

- `PORT=8080`
- `DATABASE_URL=postgres://notes:notes@localhost:5432/notes_stack?sslmode=disable`

See [env/api.env.example](/Users/mertmutlu/Documents/GitHub/notes-stack/env/api.env.example).

### 4. Verify the API

Health checks:

```bash
curl localhost:8080/healthz
curl localhost:8080/readyz
```

List notes:

```bash
curl localhost:8080/api/v1/notes
```

Create a note:

```bash
curl -X POST localhost:8080/api/v1/notes \
  -H 'Content-Type: application/json' \
  -d '{"title":"first note","content":"hello"}'
```

Get a note by ID:

```bash
curl localhost:8080/api/v1/notes/<note-id>
```

Update a note:

```bash
curl -X PUT localhost:8080/api/v1/notes/<note-id> \
  -H 'Content-Type: application/json' \
  -d '{"title":"updated title","content":"updated content"}'
```

Delete a note:

```bash
curl -i -X DELETE localhost:8080/api/v1/notes/<note-id>
```

### 5. Stop PostgreSQL

```bash
task postgres:down
```
