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
// configuration parameters can be set dynamically, like where to place the build directory
//

var root = mustRoot(magnet.Config{
	PrintConfig: true,
	LogDir:      "_build/logs",
	CacheDir:    "_build",
})

// Deinit schedules the clean up tasks to run when mage exits
var Deinit = Shutdown

var (
	goVersion = root.E(magnet.EnvVar{
		Key:     "GOLANG_VER",
		Default: "1.13.12-stretch",
		Short:   "Set the golang version to embed within the container",
	})

	golangciVersion = root.E(magnet.EnvVar{
		Key:     "GOLANGCI_VER",
		Default: "v1.27.0",
		Short:   "Set the golangci-lint version to embed within the container",
	})
)

// Build is an example mage target that builds a Dockerfile located in the present directory
func Build() (err error) {
	t := root.Target("build")
	defer func() { t.Complete(err) }()

	err = t.DockerBuild().
		SetBuildArg("GOLANG_VER", goVersion). // Passes --build-arg to docker build command
		SetBuildArg("GOLANGCI_VER", golangciVersion).
		SetPull(true).             // Same as passing --pull to the build command
		AddTag(containerName()).   // Indicate tags for the built image, can be passed multiple times
		Build(context.TODO(), ".") // Execute the docker command, passing context.Context and a docker context path of "."
	if err != nil {
		return trace.Wrap(err)
	}

	_, err = t.Exec().Run(context.TODO(), "docker", "images", "magnet-example")

	return
}

func containerName() string {
	return fmt.Sprint("magnet-example:", magnet.DefaultVersion())
}

// Run is an example mage target that runs a container
func Run() (err error) {
	t := root.Target("run")
	defer func() { t.Complete(err) }()

	// Use deps to ensure we've built the container we want to run
	mg.Deps(Build)

	// Find the WD so we can mount the PWD to the container
	wd, _ := os.Getwd()

	// Run a docker container
	err = t.DockerRun().
		SetRemove(true).                 // --rm
		SetUID(fmt.Sprint(os.Getuid())). // UID to use within the container
		SetGID(fmt.Sprint(os.Getgid())). // GID to use within the container
		AddVolume(magnet.DockerBindMount{
			Source:      wd,
			Destination: "/wd",
			Readonly:    true,
			Consistency: "cached",
		}).
		SetEnv("hello", "world").            // --env
		Run(context.TODO(), containerName(), // The imageref and command to run
			"bash",
			"-c",
			"env && ls -l /wd",
		)
	if err != nil {
		return trace.Wrap(err)
	}

	return
}

// Exec is an example mage target that execs into a running container
func Exec() (err error) {
	return trace.NotImplemented("Exec hasn't been implemented yet")
}

// Context demonstrates using a custom context to minimize the amount of data copied to the docker daemon
// This is mainly useful when working with large code bases or build directories that we want to filter out.
// Unlike the built in docker context support, the magnet support is way more flexible.
func Context() (err error) {
	t := root.Target("context")
	defer func() { t.Complete(err) }()

	err = t.DockerBuild().
		SetBuildArg("GOLANG_VER", goVersion). // Passes --build-arg to docker build command
		SetBuildArg("GOLANGCI_VER", golangciVersion).
		SetPull(true).                                    // Same as passing --pull to the build command
		AddTag("context:example").                        // Indicate tags for the built image, can be passed multiple times
		CopyToContext("assets/", "/scripts", nil, nil).   // Copy the specific file script1.sh to be /script1.sh within the docker context
		CopyToContext("build.go", "/build.go", nil, nil). // CopyToContext can be added multiple times, to selectively create the context with only needed files
		Build(context.TODO(), ".")                        // The context path here is now used only to locate the Dockerfile. Alternatively use SetDockerfile option.
	if err != nil {
		return trace.Wrap(err)
	}

	_, err = t.Exec().Run(context.TODO(), "docker", "images", "magnet-example")

	return
}

// Shutdown executes magnet's clean up tasks (internal)
func Shutdown() {
	root.Shutdown()
}

func mustRoot(config magnet.Config) *magnet.Magnet {
	root, err := magnet.Root(config)
	if err != nil {
		panic(err.Error())
	}
	return root
}
