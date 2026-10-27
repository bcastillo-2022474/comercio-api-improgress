//go:build mage

package main

import (
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

// Dev runs the server with live-reload (air).
func Dev() error {
	return sh.RunWithV(dotEnv(), "air")
}

// Run runs the server without live-reload.
func Run() error {
	return sh.RunWithV(dotEnv(), "go", "run", "./cmd/server")
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
	if err := sh.RunV("go", "build", "./..."); err != nil {
		return err
	}
	return Test()
}

type fmtError string

func (e fmtError) Error() string { return "gofmt needed on:\n" + string(e) }
