// Package container is the composition root: the one place where dependencies
// are wired together. Handlers receive already-wired use cases from here and
// never construct a repository or adapter themselves.
package container

import (
	"database/sql"
	"log/slog"
	"time"

	"joao-api-comercio/internal/config"
	dimensionsreadmodel "joao-api-comercio/internal/modules/dimensions/readmodel"
	dimensionsquery "joao-api-comercio/internal/modules/dimensions/use_cases/query"
	"joao-api-comercio/internal/modules/identity/token"
	"joao-api-comercio/internal/modules/identity/use_cases/command"
	"joao-api-comercio/internal/modules/identity/user"
	bcrypthasher "joao-api-comercio/internal/modules/identity/user/hasher"
	userrepo "joao-api-comercio/internal/modules/identity/user/repository"
	"joao-api-comercio/internal/shared/sqlite"
)

// Dependencies holds the shared infrastructure every use case is built from.
type Dependencies struct {
	Config config.Config
	Log    *slog.Logger
	DB     *sql.DB
}

// Container exposes the wired use cases.
type Container struct {
	Register *command.RegisterUseCase
	Login    *command.LoginUseCase
	Tokens   *token.JWTIssuer

	ListCountries   *dimensionsquery.ListCountriesUseCase
	ListCustoms     *dimensionsquery.ListCustomsUseCase
	ListRegimeTypes *dimensionsquery.ListRegimeTypesUseCase
}

func New(deps *Dependencies) (*Container, error) {
	tokens, err := token.NewJWTIssuer(deps.Config.JWTSecret, deps.Config.JWTTTL, time.Now)
	if err != nil {
		return nil, err
	}

	uow := sqlite.NewUnitOfWork(deps.DB, func(tx *sql.Tx) user.Repository {
		return userrepo.NewSQLiteRepository(tx)
	}, time.Now)
	hasher := bcrypthasher.NewBcrypt()

	dimensions := dimensionsreadmodel.NewSQLite(deps.DB)

	return &Container{
		Register:        command.NewRegisterUseCase(uow, hasher, time.Now, deps.Log),
		Login:           command.NewLoginUseCase(uow, hasher, tokens, deps.Log),
		Tokens:          tokens,
		ListCountries:   dimensionsquery.NewListCountriesUseCase(dimensions, deps.Log),
		ListCustoms:     dimensionsquery.NewListCustomsUseCase(dimensions, deps.Log),
		ListRegimeTypes: dimensionsquery.NewListRegimeTypesUseCase(dimensions, deps.Log),
	}, nil
}
