//go:build mage

package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/sh"
)

// dotEnv reads .env. Variables already set in the shell win.
func dotEnv() map[string]string {
	env := map[string]string{}
	raw, err := os.ReadFile(".env")
	if err != nil {
		return env
	}
	for _, line := range strings.Split(string(raw), "\n") {
		name, value, found := strings.Cut(strings.TrimSpace(line), "=")
		if !found || strings.HasPrefix(name, "#") || name == "" {
			continue
		}
		if _, set := os.LookupEnv(name); !set {
			env[name] = strings.Trim(value, `"'`)
		}
	}
	return env
}

// gooseCLI is pinned to the version the Dockerfile installs, so migrations behave the same everywhere.
const gooseCLI = "github.com/pressly/goose/v3/cmd/goose@v3.29.0"

// databasePath is DATABASE_PATH from the shell, then .env, then the default the server uses.
func databasePath() string {
	if path := os.Getenv("DATABASE_PATH"); path != "" {
		return path
	}
	if path := dotEnv()["DATABASE_PATH"]; path != "" {
		return path
	}
	return "data/comercio_exterior.db"
}

// Dev runs the server with live-reload (air).
func Dev() error {
	return sh.RunWithV(dotEnv(), "air")
}

// Run runs the server without live-reload.
func Run() error {
	return sh.RunWithV(dotEnv(), "go", "run", "./cmd/server")
}

// Setup prepares a fresh checkout: creates .env from .env.example with a new JWT_SECRET if there is
// no .env yet, then applies the migrations. Safe to run again.
func Setup() error {
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		example, err := os.ReadFile(".env.example")
		if err != nil {
			return fmt.Errorf("reading .env.example: %w", err)
		}
		secret, err := randomHex(32)
		if err != nil {
			return err
		}
		lines := strings.Split(string(example), "\n")
		for i, line := range lines {
			if strings.HasPrefix(line, "JWT_SECRET=") {
				lines[i] = "JWT_SECRET=" + secret
			}
		}
		if err := os.WriteFile(".env", []byte(strings.Join(lines, "\n")), 0o600); err != nil {
			return fmt.Errorf("writing .env: %w", err)
		}
		fmt.Println("created .env with a new JWT_SECRET")
	}
	return Migrate()
}

func randomHex(bytes int) (string, error) {
	buf := make([]byte, bytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generating JWT_SECRET: %w", err)
	}
	return hex.EncodeToString(buf), nil
}

// Build compiles the server into bin/server.
func Build() error {
	return sh.RunV("go", "build", "-o", "bin/server", "./cmd/server")
}

// Migrate applies pending database migrations to DATABASE_PATH. The server never migrates on startup.
func Migrate() error {
	return sh.RunV("go", "run", gooseCLI, "-dir", "migrations", "sqlite3", databasePath(), "up")
}

// DockerBuild builds the image locally. It is not pushed anywhere.
func DockerBuild() error {
	return sh.RunV("docker", "build", "-t", "joao-api-comercio:dev", ".")
}

// Sqlc regenerates the query packages from the *.sql files.
func Sqlc() error {
	return sh.RunV("sqlc", "generate")
}

// Test runs the test suite.
func Test() error {
	return sh.RunV("go", "test", "-count=1", "./...")
}

// Check runs everything CI would: format check, vet, build, tests.
func Check() error {
	if out, err := sh.Output("gofmt", "-l", "."); err != nil {
		return err
	} else if out != "" {
		return fmtError(out)
	}
	if err := sh.RunV("go", "vet", "./..."); err != nil {
		return err
	}
	// The magefile has its own build tag, so the vet above skips it.
	if err := sh.RunV("go", "vet", "-tags", "mage", "./..."); err != nil {
		return err
	}
	if err := sh.RunV("go", "build", "./..."); err != nil {
		return err
	}
	return Test()
}

type fmtError string

func (e fmtError) Error() string { return "gofmt needed on:\n" + string(e) }
