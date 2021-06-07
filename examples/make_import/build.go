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
	"fmt"
	"os"

	"github.com/gravitational/magnet"
	"github.com/gravitational/magnet/common"
	"github.com/magefile/mage/mg"
)

var root = mustRoot(magnet.Config{
	PrintConfig: true,
	LogDir:      "_build/logs",
	CacheDir:    "_build",
	ModulePath:  "github.com/gravitational/magnet/examples/make_import",
	ImportEnv:   mustEnv(magnet.ImportEnvFromMakefile()),
})

var (
	// vars imported from make
	goVersion = root.E(magnet.EnvVar{
		Key:   "GO_VERSION",
		Short: "Set the Go version (Default from make)",
	})

	arch = root.E(magnet.EnvVar{
		Key:   "ARCH",
		Short: "Set the arch (Default from make)",
	})

	k8sVersion = root.E(magnet.EnvVar{
		Key:   "K8S_VERSION",
		Short: "Set the k8s version (Default from make)",
	})
)

// Env runs a simple build target that imports some configuration from Make, and uses it as
// defaults for env variables.
func Env() (err error) {
	fmt.Println("arch: ", arch)
	fmt.Println("goVersion: ", goVersion)
	fmt.Println("k8sVersion: ", k8sVersion)

	return
}

type Help mg.Namespace

// Envs outputs the current environment configuration
func (Help) Envs() (err error) {
	m := root.Target("help:envs")
	defer func() { m.Complete(err) }()

	return common.WriteEnvs(root.Env(), os.Stdout)
}

func mustEnv(env map[string]string, err error) map[string]string {
	if err != nil {
		panic(err.Error())
	}
	return env
}

func mustRoot(config magnet.Config) *magnet.Magnet {
	root, err := magnet.Root(config)
	if err != nil {
		panic(err.Error())
	}
	return root
}
