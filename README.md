# joao-api-comercio

Backend API for the foreign trade dashboard: import declarations (pólizas), lookup tables, and
dashboard KPIs, served over HTTP from the provided SQLite database.

## Quick start

You need [mage](https://magefile.org) and Go 1.26. [air](https://github.com/air-verse/air) is
needed for live reload, [sqlc](https://sqlc.dev) only if you change the SQL, and Docker only to
build the image.

```sh
go install github.com/magefile/mage@latest
go install github.com/air-verse/air@latest   # only for mage dev

mage setup     # creates .env with a random JWT_SECRET, then applies the migrations
mage dev       # server on HTTP_ADDR (default :8080), reloads on save
```

Then open <http://localhost:8080/docs> for the interactive API docs. The request collection is in
[`docs/api.http`](docs/api.http) for the JetBrains or VS Code REST Client. Run the `login` request
first; the others read its token.

All the targets:

| Target | What it does |
|---|---|
| `mage setup` | creates `.env` if missing (new `JWT_SECRET`), applies migrations. Safe to rerun |
| `mage dev` | runs the server with live reload (air) |
| `mage run` | runs the server once, without live reload |
| `mage build` | compiles the server into `bin/server` |
| `mage migrate` | applies pending migrations to `DATABASE_PATH` (the server never migrates by itself) |
| `mage test` | runs the test suite |
| `mage check` | gofmt, vet, build and tests, the same checks CI runs |
| `mage sqlc` | regenerates the query code after a change to a `*.sql` file |
| `mage dockerBuild` | builds the image locally as `joao-api-comercio:dev` |

Run `mage -l` to see the list from the command line.

## Docker

```sh
mage dockerBuild
```

The image has the server and the goose CLI, not the database. To migrate and serve a database
from the host, mount its directory. On SELinux hosts (Fedora, RHEL) the mount needs `:z`:

```sh
docker run --rm -v "$PWD/data:/data:z" --entrypoint /app/goose joao-api-comercio:dev \
  -dir /app/migrations sqlite3 /data/comercio_exterior.db up

docker run --rm -p 8080:8080 -v "$PWD/data:/data:z" -e JWT_SECRET=<your secret> \
  joao-api-comercio:dev
```

The container runs as a non-root user, so the data directory has to be writable by it.

## Configuration

| Variable | Default | Meaning |
|---|---|---|
| `HTTP_ADDR` | `0.0.0.0:8080` | address the server listens on |
| `DATABASE_PATH` | `data/comercio_exterior.db` | SQLite database file |
| `JWT_SECRET` | none, required | HS256 signing key, at least 32 bytes |
| `JWT_TTL` | `15m` | access token lifetime |
| `KPI_CACHE_TTL` | `5m` | how long dashboard KPI results are cached |

The server does not migrate on startup. Run the goose command above whenever the migrations
change.

## Endpoints

Every endpoint except `/healthz`, `/docs`, `/openapi.json`, `/metrics`, `/register` and `/login`
needs `Authorization: Bearer <token>`.

| Method | Path | Purpose |
|---|---|---|
| GET | `/healthz` | liveness |
| POST | `/register` | create an account |
| POST | `/login` | get an access token |
| GET | `/dimensions/paises` | countries, for filters |
| GET | `/dimensions/aduanas` | customs offices, for filters |
| GET | `/dimensions/tipos-regimen` | customs regime types, for filters |
| POST | `/polizas` | create a declaration with its items (one transaction) |
| GET | `/polizas` | list declarations: `page`, `page_size` (max 100), `aduana_id`, `tipo_regimen_id`, `fecha_desde`, `fecha_hasta` |
| GET | `/polizas/{id}` | one declaration with its items |
| GET | `/dashboard/kpis/totals` | total value (USD), declaration count, item count |
| GET | `/dashboard/kpis/top-5` | top 5 countries and top 5 products by value |
| GET | `/dashboard/kpis/valor-por-mes` | value per month for `year` (defaults to the latest year in the data) |
| GET | `/metrics` | Prometheus metrics, meant for the internal network |

### Errors

Errors use RFC 9457 problem details (`application/problem+json`). The `type` field identifies the
error, so clients branch on it instead of parsing the message:

```json
{
  "type": "urn:joao-api-comercio:error:invalid_fecha_declaracion",
  "title": "Unprocessable Entity",
  "status": 422,
  "detail": "fecha_declaracion must be YYYY-MM-DD, got \"2025-01-011\""
}
```

Auth failures use `missing_token` and `invalid_token`. Unexpected failures return a generic 500
with no internal detail; the cause is logged once, at the HTTP edge.

## Architecture

The code is organized as bounded contexts under `internal/modules`, each split into layers:

- **Domain** (`poliza`, `user`, `kpi`): aggregates, invariants, domain errors and their codes.
  No SQL, no HTTP.
- **Ports** (`ports.go`): repository interfaces for writes. Each has an in-memory fake and a
  contract suite that the fake and the SQL adapter both must pass.
- **Read models** (`readmodel`): screen-shaped DTOs and one query interface per query. Reads
  never go through repositories.
- **Use cases** (`use_cases/command`, `use_cases/query`): one file per use case. Every use case
  logs its start and its outcome at Info. Errors are returned, not logged.
- **Handlers** (`handler`): Huma operations. They translate HTTP to use cases and domain errors
  to HTTP through one mapping (`internal/shared/httperr`).
- **Composition root** (`internal/container`): the only place that wires adapters to use cases.

Decisions:

- **Huma** for the HTTP layer. Routes, validation and the OpenAPI spec come from the same types.
  Query parameters that may be absent use `httpx.OptionalParam`, because Huma does not allow
  pointers in parameters.
- **SQLite with plain SQL.** Queries are written by hand and compiled by sqlc, so every query is
  checked against the schema at generation time. sqlc column overrides map nullable columns to
  pointers and the DATE columns to text. The read-side filters use `sqlc.narg`.
- **Migrations with goose**, embedded in the binary. Each migration's down step is written, except
  for the import tables, which hold data we did not create.
- **Writes go through a unit of work.** A command writes its rows and its domain events in one SQL
  transaction. Events go to an outbox table. The outbox relay is not built yet, so nothing
  publishes those events.
- **Passwords** are hashed with bcrypt. Login returns the same error for an unknown email and a wrong
  password. Access tokens are HS256 JWTs with a short lifetime.
- **Errors and logs.** Each domain error has a code, and a single map turns codes into HTTP
  statuses. A test checks that every declared code has a status, so a new code cannot ship
  without one.

## Testing

```sh
mage test     # go test -count=1 ./...
mage check    # gofmt, vet, build and tests
```

The suite covers the domain (entity rules), the use cases (on in-memory fakes), the repository
contracts (run against the fake and against SQLite), the error mapping, and the TTL cache (expiry,
overwrite, concurrent access under `-race`). There are no tests for the JWT and password adapters
and no HTTP-level tests.

## Observability

- Structured logs (`log/slog`). Each use case adds a `use_case` attribute. Failed logins are
  logged as warnings, without the email address.
- `GET /metrics` exposes request counts and durations per route pattern, plus the Go runtime and
  process collectors.

## Data

The API reads a database produced by the assignment's `setup.py`, so some things come from that
data:

- `setup.py` declares its `NOT NULL` constraints as `NOT_NULL`, which SQLite does not enforce. Our
  migrations use the correct constraint on fresh databases.
- 2,262 declarations have no customs office or regime. Those fields are nullable in the API.
- Some dates are from 1970. They are returned as they are in the source.

## Not done yet

- Docker Compose: there is a Dockerfile, but no compose file.
- Index migration for the declarations list. The list is a full table scan today.
- Refresh tokens and logout. Access tokens are valid until they expire.
