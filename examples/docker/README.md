# Docker

## Build
Build a docker container

`go run mage.go build`

```
❯ go run mage.go build
Logs:    _build/logs/latest (_build/logs/20200723042630)
Version: v0.1.0-17-g29f9b23-dirty
Cache:   _build/magnet/github.com/gravitational/magnet/examples/docker
[+] Building 2.2s (1/1) FINISHED
 => build                                                                                                                               2.2s

❯ cat _build/logs/latest/build
Name: build
Digest: sha256:44575cf5b28512d75644bf54a517dcef304ff809fd511747621b4d64f19aac66
Cached: false
Started: 2020-07-23 04:26:30.730309484 +0000 UTC m=+0.010997756
Completed: <nil>
-----
Env:  map[DOCKER_BUILDKIT:1 PROGRESS_NO_TRUNC:1]  Exec:  docker build --pull --compress --build-arg GOLANGCI_VER=v1.27.0 --build-arg GOLANG_VER=1.13.12-stretch -t magnet-example:v0.1.0-17-g29f9b23-dirty .
#1 [internal] load .dockerignore
#1 sha256:ff8a320c49e988bcc0dc50f5e2ee53e4d993018a51a7fbb874c578ac99e1077c
#1 transferring context: 2B 0.1s done
#1 DONE 0.1s

#2 [internal] load build definition from Dockerfile
#2 sha256:718f64813e281739da7006be1bb11e99b55fc255778a1d7b3eb1bdcf8be81f3a
#2 transferring dockerfile: 38B 0.1s done
#2 DONE 0.2s

#3 resolve image config for docker.io/docker/dockerfile:1.1.7-experimental
#3 sha256:adf0ad58625e45ab290e1bebbcfe1d73aceb4d8dbec09db8a7422cf74eff2996
#3 DONE 0.5s

#4 docker-image://docker.io/docker/dockerfile:1.1.7-experimental@sha256:de85b2f3a3e8a2f7fe48e8e84a65f6fdd5cd5183afa6412fff9caa6871649c44
#4 sha256:4f557cac37259deb5e7fe213b122cca6da5945b2c1d8b3b3f23f2bb89e44211c
#4 CACHED

#5 [internal] load metadata for quay.io/gravitational/debian-venti:go1.13.12-stretch
#5 sha256:909bba4807a2260a72db645a7249e9fdcb3e74b3d1d8c08a986dc9d6b53cf48c
#5 DONE 0.2s

#8 [stage-0 1/3] FROM quay.io/gravitational/debian-venti:go1.13.12-stretch@sha256:2572b4f2661ebbf546046e271b309079b1e9444027fb7f27be2a27cb95343dfa
#8 sha256:6cea9c97a56a62d3ac18b8b6f9a9dc4a03475fa9b325826d326ecad7a4b9442b
#8 DONE 0.0s

#6 [stage-0 2/3] RUN curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b /gopath/bin v1.27.0
#6 sha256:1acd0d5f69382081a1f3ca1f076afbd7c9df40400390e052aa7fa85bd5fcbf60
#6 CACHED

#7 [stage-0 3/3] RUN --mount=type=cache,sharing=locked,id=aptlib,target=/var/lib/apt     --mount=type=cache,sharing=locked,id=aptcache,target=/var/cache/apt           DEBIAN_FRONTEND=noninteractive apt-get -q -y update --fix-missing && apt-get -q -y install libaio-dev zlib1g-dev
#7 sha256:acb06cdbfa09e9e34184d9333312af5b28a0221b7123a63128c97b64d87e9ee7
#7 CACHED

#9 exporting to image
#9 sha256:e8c613e07b0b7ff33893b694f7759a10d42e180f2b4dc349fb57dc6b71dcab00
#9 exporting layers done
#9 writing image sha256:846ecfff0f87e3c60fae109060b387a92c34c4a9d8f62046d095eda7486e6de8 done
#9 naming to docker.io/library/magnet-example:v0.1.0-17-g29f9b23-dirty done
#9 DONE 0.0s
Exec:  docker images magnet-example
REPOSITORY          TAG                        IMAGE ID            CREATED             SIZE
magnet-example      v0.1.0-17-g29f9b23-dirty   846ecfff0f87        2 days ago          1.16GB
-----
Vertex: Completed <nil> -> 2020-07-23 04:26:32.887329557 +0000 UTC m=+2.168017861
Vertex: Duration 2.157020105s
-----


```

## Run
Runs a command within a container, with a dependency on the build target.

`go run mage.go run`

```
❯ go run mage.go run
Logs:    _build/logs/latest (_build/logs/20200723042741)
Version: v0.1.0-17-g29f9b23-dirty
Cache:   _build/magnet/github.com/gravitational/magnet/examples/docker
[+] Building 2.3s (2/2) FINISHED
 => run                                                                                                                                 2.3s
 => build                                                                                                                               1.0s

❯ cat _build/logs/latest/run
Name: run
Digest: sha256:acba25512100f80b56fc3ccd14c65be55d94800cda77585c5f41a887e398f9be
Cached: false
Started: 2020-07-23 04:27:41.76455901 +0000 UTC m=+0.009779786
Completed: <nil>
-----
Exec:  docker run -u 1003:1005 --rm=true --mount type=bind,source=/home/knisbet/go/src/github.com/gravitational/magnet/examples/docker,target=/wd,readonly,consistency=cached --env=hello=world magnet-example:v0.1.0-17-g29f9b23-dirty bash -c env && ls -l /wd
HOSTNAME=e959394b999b
GOPATH=/gopath
PWD=/
HOME=/
GOROOT=/go
DEBIAN_FRONTEND=noninteractive
SHLVL=1
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/go/bin
hello=world
_=/usr/bin/env
total 24
-rw-rw-r-- 1 1003 1005  630 Jul 21 04:43 Dockerfile
-rw-rw-r-- 1 1003 1005 1768 Jul 23 04:26 README.md
drwxrwxr-x 2 1003 1005 4096 Jul 22 04:18 assets
drwxr-xr-x 1 1003 1005   60 Jul 21 04:54 build
-rw-rw-r-- 1 1003 1005 4570 Jul 22 15:35 build.go
-rw-rw-r-- 1 1003 1005  125 Jul 21 04:38 mage.go
-----
Vertex: Completed <nil> -> 2020-07-23 04:27:44.019278784 +0000 UTC m=+2.264499569
Vertex: Duration 2.254719783s
-----

```


## Context
Build a docker container, but use a narrow context that only includes needed files.

`go run mage.go context`

```
❯ go run mage.go context
Logs:    _build/logs/latest (_build/logs/20200723042904)
Version: v0.1.0-17-g29f9b23-dirty
Cache:   _build/magnet/github.com/gravitational/magnet/examples/docker
[+] Building 1.3s (1/1) FINISHED
 => context                                                                                                                             1.3s

❯ cat _build/logs/latest/context
Name: context
Digest: sha256:ea7792a26f405e2ae9c6f49ca93bbe6076ceac0a1fc53d83426c7d7f2d9377e4
Cached: false
Started: 2020-07-23 04:29:04.983942225 +0000 UTC m=+0.010395866
Completed: <nil>
-----
Env:  map[DOCKER_BUILDKIT:1 PROGRESS_NO_TRUNC:1]  Exec:  docker build --pull --compress --build-arg GOLANGCI_VER=v1.27.0 --build-arg GOLANG_VER=1.13.12-stretch -t context:example /tmp/docker-context249320160
#2 [internal] load .dockerignore
#2 sha256:b2f26ebc913dd3436f504e7604b5b9feef0e0e4f360414d5fe42529e7acb648a
#2 transferring context: 2B done
#2 DONE 0.1s

#1 [internal] load build definition from Dockerfile
#1 sha256:1033a000f6b40ffe0e79bf81573adf480de84c7b763317b15479248866df8410
#1 transferring dockerfile: 675B done
#1 DONE 0.1s

#3 resolve image config for docker.io/docker/dockerfile:1.1.7-experimental
#3 sha256:adf0ad58625e45ab290e1bebbcfe1d73aceb4d8dbec09db8a7422cf74eff2996
#3 DONE 0.3s

#4 docker-image://docker.io/docker/dockerfile:1.1.7-experimental@sha256:de85b2f3a3e8a2f7fe48e8e84a65f6fdd5cd5183afa6412fff9caa6871649c44
#4 sha256:4f557cac37259deb5e7fe213b122cca6da5945b2c1d8b3b3f23f2bb89e44211c
#4 CACHED

#5 [internal] load metadata for quay.io/gravitational/debian-venti:go1.13.12-stretch
#5 sha256:909bba4807a2260a72db645a7249e9fdcb3e74b3d1d8c08a986dc9d6b53cf48c
#5 DONE 0.1s

#6 [stage-0 1/3] FROM quay.io/gravitational/debian-venti:go1.13.12-stretch@sha256:2572b4f2661ebbf546046e271b309079b1e9444027fb7f27be2a27cb95343dfa
#6 sha256:6cea9c97a56a62d3ac18b8b6f9a9dc4a03475fa9b325826d326ecad7a4b9442b
#6 DONE 0.0s

#7 [stage-0 2/3] RUN curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b /gopath/bin v1.27.0
#7 sha256:1acd0d5f69382081a1f3ca1f076afbd7c9df40400390e052aa7fa85bd5fcbf60
#7 CACHED

#8 [stage-0 3/3] RUN --mount=type=cache,sharing=locked,id=aptlib,target=/var/lib/apt     --mount=type=cache,sharing=locked,id=aptcache,target=/var/cache/apt           DEBIAN_FRONTEND=noninteractive apt-get -q -y update --fix-missing && apt-get -q -y install libaio-dev zlib1g-dev
#8 sha256:acb06cdbfa09e9e34184d9333312af5b28a0221b7123a63128c97b64d87e9ee7
#8 CACHED

#9 exporting to image
#9 sha256:e8c613e07b0b7ff33893b694f7759a10d42e180f2b4dc349fb57dc6b71dcab00
#9 exporting layers done
#9 writing image sha256:846ecfff0f87e3c60fae109060b387a92c34c4a9d8f62046d095eda7486e6de8 done
#9 naming to docker.io/library/context:example done
#9 DONE 0.0s
Exec:  docker images magnet-example
REPOSITORY          TAG                        IMAGE ID            CREATED             SIZE
magnet-example      v0.1.0-17-g29f9b23-dirty   846ecfff0f87        2 days ago          1.16GB
-----
Vertex: Completed <nil> -> 2020-07-23 04:29:06.255924245 +0000 UTC m=+1.282377880
Vertex: Duration 1.271982014s

```
