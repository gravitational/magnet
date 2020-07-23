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
	"time"

	"github.com/gravitational/magnet"

	// mage:import
	_ "github.com/gravitational/magnet/common"
)

var root = magnet.Root(magnet.Config{
	PrintConfig: true,
})

// HelloWorld runs a simple build target that just prints some output and exits
func HelloWorld() (err error) {
	m := root.Target("helloworld")
	defer func() { m.Complete(err) }()

	for i := 1; i <= 10; i++ {
		m.Println("hello world: ", i)
		time.Sleep(500 * time.Millisecond)
	}

	return
}
