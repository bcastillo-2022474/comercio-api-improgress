# Build stage: a static server binary, plus the goose CLI for running migrations.
FROM golang:1.26-bookworm AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server
RUN CGO_ENABLED=0 go install -trimpath -ldflags="-s -w" github.com/pressly/goose/v3/cmd/goose@v3.29.0

# Runtime stage: both binaries are static, so no libc or extra packages are needed.
FROM alpine:3
RUN adduser -D -u 10001 app
WORKDIR /app

COPY --from=build /out/server /app/server
COPY --from=build /go/bin/goose /app/goose
COPY migrations /app/migrations

RUN mkdir /data && chown app /data
USER app

ENV HTTP_ADDR=0.0.0.0:8080 \
    DATABASE_PATH=/data/comercio_exterior.db

EXPOSE 8080
ENTRYPOINT ["/app/server"]
