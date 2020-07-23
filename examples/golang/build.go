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

	"github.com/gravitational/magnet"
	"github.com/gravitational/trace"

	// mage:import
	_ "github.com/gravitational/magnet/common"
)

//
// configuration parameters can be set dynamically, like where to place the build directory
//

var root = magnet.Root(magnet.Config{
	PrintConfig: true,
})

var (
	goVersion = magnet.E(magnet.EnvVar{
		Key:     "GOLANG_VER",
		Default: "1.13.12-stretch",
		Short:   "Set the golang version to embed within the container",
	})
)

// Build is an example mage target that builds a golang project
// Note: currently assumes project extracted to a GOPATH directory structure
func Build() (err error) {
	t := root.Target("build")
	defer func() { t.Complete(err) }()

	err = t.GolangBuild().
		SetOutputPath("build/example.local").
		Build(context.TODO(), "github.com/gravitational/magnet/examples/golang")
	if err != nil {
		return trace.Wrap(err)
	}
	return
}

// BuildContainer is an example mage target that builds a golang project within a docker container
// Note: currently assumes project extracted to a GOPATH directory structure
func BuildContainer() (err error) {
	t := root.Target("build")
	defer func() { t.Complete(err) }()

	err = t.GolangBuild().
		SetOutputPath("build/example.container").
		SetBuildContainer("golang:1.14").
		Build(context.TODO(), "github.com/gravitational/magnet/examples/golang")
	if err != nil {
		return trace.Wrap(err)
	}
	return
}

// Test is an example mage target that runs the go compilers tests
// Note: currently assumes project extracted to a GOPATH directory structure
func Test() (err error) {
	t := root.Target("test")
	defer func() { t.Complete(err) }()

	err = t.GolangTest().
		SetBuildContainer("golang:1.14").
		Test(context.TODO(), "github.com/gravitational/magnet/examples/golang")
	if err != nil {
		return trace.Wrap(err)
	}
	return
}
