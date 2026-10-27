//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Run starts the API server.
func Run() error {
	return sh.RunV("go", "run", "./cmd/server")
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
