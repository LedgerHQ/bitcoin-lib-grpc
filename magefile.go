// +build mage

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/magefile/mage/sh"
)

const (
	entryPoint = "cmd/lbs.go"
	ldFlags    = "-X $PACKAGE/version/version.commitHash=$COMMIT_HASH " +
		"-X $PACKAGE/version/version.buildDate=$BUILD_DATE"
	protoPlugins  = "plugins=grpc"
	protoDir      = "pb/bitcoin"
	protoFileName = "service.proto"
)

// Allow user to override executables on UNIX-like systems.
var goexe = "go"      // GOEXE=xxx mage build
var protoc = "protoc" // PROTOC=xxx mage proto
var buf = "buf"       // BUF=xxx mage protolist

func init() {
	if exe := os.Getenv("GOEXE"); exe != "" {
		goexe = exe
	}

	if exe := os.Getenv("PROTOC"); exe != "" {
		protoc = exe
	}

	if exe := os.Getenv("BUF"); exe != "" {
		buf = exe
	}

	// We want to use Go 1.11 modules even if the source lives inside GOPATH.
	// The default is "auto".
	os.Setenv("GO111MODULE", "on")
}

func Proto() error {
	return sh.Run(protoc,
		fmt.Sprintf("--go_out=%s:%s", protoPlugins, protoDir), // protoc flags
		fmt.Sprintf("%s/%s", protoDir, protoFileName))         // input .proto
}

func Buf() error {
	// Verify if the proto files can be compiled.
	if err := sh.Run(buf, "image", "build", "-o /dev/null"); err != nil {
		return err
	}

	// Run Buf lint checks on the protobuf files.
	if err := sh.Run(buf, "check", "lint"); err != nil {
		return err
	}

	return nil
}

// Build binary
func Build() error {
	if err := Proto(); err != nil {
		return err
	}

	return sh.RunWith(flagEnv(), goexe, "build", "-ldflags", ldFlags,
		entryPoint)
}

// Run tests
func Test() error {
	return sh.Run(goexe, "test", "./...")
}

// Run tests with race detector
func TestRace() error {
	return sh.Run(goexe, "test", "-race", "./...")
}

// Run tests with race-detector and code-coverage.
// Useful on CI, but can be run locally too.
func TestRaceCover() error {
	return sh.Run(
		goexe, "test", "-race", "-coverprofile=coverage.txt",
		"-covermode=atomic", "./...")
}

// Run basic golangci-lint check.
func Lint() error {
	linterArgs := []string{
		"run",
		"--disable-all",
		"--enable=govet",
		"--enable=gofmt",
		"--enable=gosec",
	}

	if err := sh.Run("golangci-lint", linterArgs...); err != nil {
		return err
	}

	return nil
}

func flagEnv() map[string]string {
	hash, _ := sh.Output("git", "rev-parse", "--short", "HEAD")
	return map[string]string{
		"PACKAGE":     entryPoint,
		"COMMIT_HASH": hash,
		"BUILD_DATE":  time.Now().Format("2006-01-02T15:04:05Z0700"),
	}
}
