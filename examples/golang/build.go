//+build mage

/*
Copyright 2020 Gravitational, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gravitational/magnet"
	"github.com/gravitational/trace"

	"github.com/magefile/mage/mg"
)

//
// configuration parameters can be set dynamically, like where to place the logs/cache
//

var root = mustRoot(magnet.Config{
	PrintConfig: true,
	LogDir:      "_build/logs",
	CacheDir:    "_build",
	ModulePath:  "github.com/gravitational/magnet/examples/golang",
	// If VERSION is unspecified, the default will be generated based on a git tag
	Version: os.Getenv("VERSION"),
})

// Deinit schedules the clean up tasks to run when mage exits
var Deinit = Shutdown

var (
	goVersion = root.E(magnet.EnvVar{
		Key:     "GOLANG_VER",
		Default: "1.13.12-stretch",
		Short:   "Set the Go version to embed within the container",
	})
	golangciVersion = root.E(magnet.EnvVar{
		Key:     "GOLANGCI_VER",
		Default: "v1.40.1",
		Short:   "Set the golangci-lint version to embed within the container",
	})
	// version of the output / build container tag
	// Can be overridden on command line with:
	//
	// # VERSION=1.0 go run mage.go buildContainer
	version = root.E(magnet.EnvVar{
		Key:     "VERSION",
		Default: root.Version,
		Short:   "Set the output version",
	})
)

// Build is an example mage target that builds a golang project
// Note: currently assumes project extracted to a GOPATH directory structure
func Build() (err error) {
	t := root.Target("build")
	defer func() { t.Complete(err) }()

	err = t.GolangBuild().
		SetOutputPath("_build/example.local").
		Build(context.TODO(), "github.com/gravitational/magnet/examples/golang")
	if err != nil {
		return trace.Wrap(err)
	}
	return
}

// BuildInContainer is an example mage target that builds a golang project within a docker container
// Note: currently assumes project extracted to a GOPATH directory structure
func BuildInContainer(ctx context.Context) (err error) {
	t := root.Target("buildInContainer")
	defer func() { t.Complete(err) }()

	mg.CtxDeps(ctx, BuildContainer)

	err = t.GolangBuild().
		SetOutputPath("_build/example.container").
		SetBuildContainer(buildContainer()).
		SetEnv("GO111MODULE", "on").
		Build(ctx, "github.com/gravitational/magnet/examples/golang")
	if err != nil {
		return trace.Wrap(err)
	}
	return
}

// Test is an example mage target that runs the go compilers tests
// Note: currently assumes project extracted to a GOPATH directory structure
func Test(ctx context.Context) (err error) {
	t := root.Target("test")
	defer func() { t.Complete(err) }()

	err = t.GolangTest().
		SetEnv("GO111MODULE", "on").
		SetBuildContainer(buildContainer()).
		Test(ctx, "github.com/gravitational/magnet/examples/golang")
	if err != nil {
		return trace.Wrap(err)
	}
	return
}

// BuildContainer creates a docker container as a consistent Go environment to use for software builds.
func BuildContainer(ctx context.Context) (err error) {
	m := root.Target("buildContainer")
	defer func() { m.Complete(err) }()

	return m.DockerBuild().
		AddTag(buildContainer()).
		SetPull(true).
		SetBuildArg("GOLANG_VER", goVersion).
		SetBuildArg("GOLANGCI_VER", golangciVersion).
		SetBuildArg("UID", fmt.Sprint(os.Getuid())).
		SetBuildArg("GID", fmt.Sprint(os.Getgid())).
		SetDockerfile("Dockerfile").
		Build(ctx, ".")
}

// Shutdown executes magnet's clean up tasks (internal)
func Shutdown() {
	root.Shutdown()
}

func buildContainer() string {
	return fmt.Sprint("build:", version)
}

func mustRoot(config magnet.Config) *magnet.Magnet {
	root, err := magnet.Root(config)
	if err != nil {
		panic(err.Error())
	}
	return root
}
