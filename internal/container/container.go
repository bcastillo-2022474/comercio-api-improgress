// Package container is the composition root: the one place where dependencies
// are wired together. Handlers receive already-wired use cases from here and
// never construct a repository or adapter themselves.
package container

import (
	"log/slog"

	"joao-api-comercio/internal/config"
)

// Dependencies holds the shared infrastructure every use case is built from.
type Dependencies struct {
	Config config.Config
	Log    *slog.Logger
}

// Container exposes the wired use cases. It is empty until the first module exists.
type Container struct {
	deps *Dependencies
}

func New(deps *Dependencies) *Container {
	return &Container{deps: deps}
}
