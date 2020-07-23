# Golang

## Build using the local compiler
`go run mage.go build`

```
❯ go run mage.go build
Logs:    build/logs/latest (build/logs/20200723041555)
Version:  v0.1.0-17-g29f9b23-dirty
Build:    build/v0.1.0-17-g29f9b23-dirty
[+] Building 0.5s (1/1) FINISHED
 => build

❯ build/example.local
Hello World

```

## Build using a container
`go run mage.go buildContainer`

```
❯ go run mage.go buildContainer
Logs:    build/logs/latest (build/logs/20200723042056)
Version:  v0.1.0-17-g29f9b23-dirty
Build:    build/v0.1.0-17-g29f9b23-dirty
[+] Building 0.8s (1/1) FINISHED
 => build

❯ build/example.container
Hello World

```


## Run tests
`go run mage.go test`

```
❯ go run mage.go test
Logs:    build/logs/latest (build/logs/20200723042430)
Version:  v0.1.0-17-g29f9b23-dirty
Build:    build/v0.1.0-17-g29f9b23-dirty
[+] Building 0.7s (1/1) FINISHED
 => test

❯ cat build/logs/latest/test
Name: test
Digest: sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08
Cached: false
Started: 2020-07-23 04:24:30.161105775 +0000 UTC m=+0.009806518
Completed: <nil>
-----
Exec:  docker run -u 1003:1005 --rm=true --mount type=bind,source=/home/knisbet/go/src/github.com/gravitational/magnet/examples/golang,target=/go/src/github.com/gravitational/magnet/examples/golang,consistency=delegated --mount type=bind,source=/home/knisbet/go/src/github.com/gravitational/magnet/examples/golang/build/cache,target=/cache,consistency=delegated --env=XDG_CACHE_HOME=/cache --env=GOCACHE=/cache/go golang:1.14 go test github.com/gravitational/magnet/examples/golang
?   	github.com/gravitational/magnet/examples/golang	[no test files]
-----
Vertex: Completed <nil> -> 2020-07-23 04:24:30.841578958 +0000 UTC m=+0.690279791
Vertex: Duration 680.473273ms

```