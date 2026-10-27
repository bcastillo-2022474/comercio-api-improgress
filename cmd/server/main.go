package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/cors"

	"joao-api-comercio/internal/config"
	"joao-api-comercio/internal/container"
	dimensionshandler "joao-api-comercio/internal/modules/dimensions/handler"
	identityhandler "joao-api-comercio/internal/modules/identity/handler"
	"joao-api-comercio/internal/shared/httpx"
	"joao-api-comercio/internal/shared/sqlite"
)

func main() {
	logger := slog.New(log.NewWithOptions(os.Stdout, log.Options{
		Level:           log.DebugLevel,
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
	}))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("loading config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := sqlite.Open(ctx, cfg.DatabasePath)
	if err != nil {
		logger.Error("opening database", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	app, err := container.New(&container.Dependencies{Config: cfg, Log: logger, DB: db})
	if err != nil {
		logger.Error("setting up dependencies", "err", err)
		os.Exit(1)
	}
	mux := http.NewServeMux()
	apiConfig := huma.DefaultConfig("joao-api-comercio", "0.1.0")
	apiConfig.DocsRenderer = huma.DocsRendererStoplightElements
	api := humago.New(mux, apiConfig)
	registerHealth(api)
	identityhandler.NewIdentity(app.Register, app.Login, logger).Mount(api)
	dimensionshandler.NewDimensions(app.ListCountries, app.ListCustoms, app.ListRegimeTypes, logger).Mount(api)

	server := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: cors.AllowAll().Handler(httpx.Logging(logger, mux)),
	}

	go func() {
		<-ctx.Done()
		_ = server.Shutdown(context.Background())
	}()

	logger.Info("server running", "addr", cfg.HTTPAddr, "docs", "http://"+cfg.HTTPAddr+"/docs")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}

type healthOutput struct {
	Body struct {
		Status string `json:"status" example:"ok"`
	}
}

// registerHealth exposes a liveness endpoint for container healthchecks.
func registerHealth(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/healthz",
		Summary:     "Liveness check",
		Tags:        []string{"system"},
	}, func(_ context.Context, _ *struct{}) (*healthOutput, error) {
		out := &healthOutput{}
		out.Body.Status = "ok"
		return out, nil
	})
}
