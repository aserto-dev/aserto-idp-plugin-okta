// +build mage

package main

import (
	"os"
	"runtime"

	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func init() {
	// Set go version for docker builds
	os.Setenv("GO_VERSION", "1.16")
	// Set private repositories
	os.Setenv("GOPRIVATE", "github.com/aserto-dev")
}

// Build builds all binaries in ./cmd.
func Build() error {
	return common.BuildReleaser()
}

// Cleans the bin director
func Clean() error {
	return os.RemoveAll("dist")
}

// BuildAll builds all binaries in ./cmd for
// all configured operating systems and architectures.
func BuildAll() error {
	return common.BuildAll()
}

func Deps() {
	deps.GetAllDeps()
}

// Lint runs linting for the entire project.
func Lint() error {
	return common.Lint()
}

// Test runs all tests and generates a code coverage report.
func Test() error {
	return common.Test()
}

func Generate() error {
	return common.Generate()
}

// All runs all targets in the appropriate order.
// The targets are run in the following order:
// deps, generate, lint, test, build, dockerImage
func All() error {
	mg.SerialDeps(Deps, Lint, Test, Build)
	return nil
}

func Run() error {
	return sh.RunV("./bin/" + runtime.GOOS + "-" + runtime.GOARCH + "/aserto-idp")
}
