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
	"time"

	"github.com/gravitational/magnet"
	"github.com/gravitational/trace"
	"github.com/magefile/mage/mg"

	// mage:import
	_ "github.com/gravitational/magnet/common"
)

//
// configuration parameters can be set dynamically, like where to place the build directory
//

var root = magnet.Root(magnet.Config{
	Version:     version,
	LogDir:      magnet.DefaultLogDir(),
	BuildDir:    magnet.DefaultBuildDir(version),
	PrintConfig: true,
})

//
// Run time parameters can be set by using E to get Environment variables, with defaults and descriptions
//

var (
	version = magnet.E(magnet.EnvVar{
		Key:     "VERSION",
		Default: magnet.DefaultVersion(),
		Short:   "Set the version that wll be built",
	})
)

// MultipleTargets is a mage build target, available by calling mage
//
// This example demonstrates creating multiple targets in a hierarchy, demonstrating how individual targets
// can be children of root or other targets, helping show the relationship between targets.
func MultipleTargets() (err error) {
	t1 := root.Target("target1")
	t1.SetCached(true)
	defer func() { t1.Complete(err) }()

	t2 := t1.Target("target2")
	defer func() { t2.Complete(err) }()

	for i := 1; i <= 1; i++ {
		t2.Println("2: hello world: ", i)
		time.Sleep(500 * time.Millisecond)
	}

	t3 := t2.Target("target3")
	defer func() { t3.Complete(err) }()

	for i := 1; i <= 2; i++ {
		t3.Println("3: hello world: ", i)
		time.Sleep(500 * time.Millisecond)
	}

	t4 := t2.Target("target4")
	defer func() { t4.Complete(err) }()

	for i := 1; i <= 2; i++ {
		t4.Println("4: hello world: ", i)
		time.Sleep(500 * time.Millisecond)
	}

	t5 := t3.Target("target5")
	var err5 error
	defer func() { t5.Complete(err5) }()

	for i := 1; i <= 2; i++ {
		t5.Println("5: hello world: ", i)
		time.Sleep(500 * time.Millisecond)
	}

	err5 = fmt.Errorf("Error on target 5")
	t5.Println("Error: ", err5.Error())

	time.Sleep(2 * time.Second)
	t1.Println("Ending")
	t2.Println("Ending")
	t3.Println("Ending")
	t4.Println("Ending")
	t5.Println("Ending")
	time.Sleep(3 * time.Second)

	return
}

func Dl() (err error) {
	t := root.Target("dl")
	defer func() { t.Complete(err) }()

	mg.Deps(Dep1, Dep2)

	var path string
	path, err = t.Download("https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl")
	//path, err = t.Download("http://ipv4.download.thinkbroadband.com/1GB.zip")

	t.Println("Path: ", path)
	return
}

// DlParallel runs multiple downloads in parallel
func DlParallel() (err error) {
	t := root.Target("downloads")
	defer func() { t.Complete(err) }()

	mg.Deps(Dep1, Dep2)

	kubectl := t.DownloadFuture("https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl")
	gb := t.DownloadFuture("http://ipv4.download.thinkbroadband.com/50MB.zip")
	bad := t.DownloadFuture("http://example.com/non-existant-file")

	var errors []error
	for _, future := range []func() (string, string, error){kubectl, gb, bad} {
		url, path, err := future()
		t.Printlnf("url: %v path: %v error: %v", url, path, trace.DebugReport(err))
		if err != nil {
			errors = append(errors, err)
		}
	}

	// Simulate some amount of work
	time.Sleep(5 * time.Second)
	return trace.NewAggregate(errors...)
}

// Dep1 is executed as a dependency of the DL tasks
func Dep1() (err error) {
	t := root.Target("dep1")
	defer func() { t.Complete(err) }()

	_, err = t.Download("https://speed.hetzner.de/100MB.bin")
	return
}

// Dep2 is executed as a dependency of the DL tasks
func Dep2() (err error) {
	t := root.Target("dep2")
	defer func() { t.Complete(err) }()

	_, err = t.Download("http://speedtest-ny.turnkeyinternet.net/100mb.bin")

	return
}
