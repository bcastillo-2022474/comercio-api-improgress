# joao-api-comercio

Go backend skeleton: HTTP API built with [Huma](https://huma.rocks) on `net/http`,
composition root in `internal/container`, one directory per bounded context under
`internal/modules`.

## Run

```
go install github.com/magefile/mage@latest
mage run      # starts the server on :8080 (override with HTTP_ADDR)
mage check    # gofmt, vet, build, tests
```

- Health: `GET /healthz`
- OpenAPI docs: `GET /docs` (spec at `/openapi.json`)

## Layout

```
cmd/server/          entry point: config, container, routes, graceful shutdown
internal/config/     env-var config
internal/container/  composition root, where dependencies are wired
internal/modules/    one directory per bounded context (empty for now)
internal/shared/     cross-cutting code: HTTP middleware
```

Not included yet: domain modules, database, migrations, message broker.
