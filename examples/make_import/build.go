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

	"github.com/gravitational/magnet"

	// mage:import
	_ "github.com/gravitational/magnet/common"
)

var root = magnet.Root(magnet.Config{
	PrintConfig:   true,
	MakeEnvTarget: "print-vars",
})

var (
	// vars imported from make
	goVersion = magnet.E(magnet.EnvVar{
		Key:   "GO_VERSION",
		Short: "Set the golang version (Default from make)",
	})

	arch = magnet.E(magnet.EnvVar{
		Key:   "ARCH",
		Short: "Set the arch (Default from make)",
	})

	k8sVersion = magnet.E(magnet.EnvVar{
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
