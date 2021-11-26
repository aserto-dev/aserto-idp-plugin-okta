// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"runtime"

	"github.com/aserto-dev/mage-loot/common"
	"github.com/aserto-dev/mage-loot/deps"
	"github.com/aserto-dev/sver/pkg/sver"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

func init() {
	// Set go version for docker builds
	os.Setenv("GO_VERSION", "1.16")
	// Set private repositories
	os.Setenv("GOPRIVATE", "github.com/aserto-dev")
}

var (
	oras         = deps.BinDep("oras")
	mediaType    = ":application/vnd.unknown.layer.v1+txt"
	pluginName   = "aserto-idp-plugin-okta_"
	execLocation = "dist/" + pluginName
	ghName       = "ghcr.io/aserto-dev/" + pluginName
	osArch       = []string{"linux_arm64", "linux_amd64", "darwin_amd64", "darwin_arm64", "windows_amd64"}
)

// Build builds all binaries in ./cmd.
func Build() error {
	return common.BuildReleaser()
}

// Cleans the bin director
func Clean() error {
	return os.RemoveAll("dist")
}

// Release releases the project.
func Release() error {
	return common.Release()
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

func Publish() error {

	username := os.Getenv("DOCKER_USERNAME")
	if username == "" {
		return errors.New("env var DOCKER_USERNAME is not set")
	}
	password := os.Getenv("DOCKER_PASSWORD")
	if password == "" {
		return errors.New("env var DOCKER_PASSWORD is not set")
	}

	version, err := sver.CurrentVersion(true, true)
	fmt.Println(version)
	if err != nil {
		return fmt.Errorf("couldn't calculate current version: %w", err)
	}

	for _, arch := range osArch {
		err = oras("push", "-u", username, "-p", password, ghName+arch+":"+version, execLocation+arch+mediaType)
		if err != nil {
			return err
		}
	}

	return nil
}
