# syntax=docker/dockerfile:1.1.7-experimental

ARG GOLANG_VER
FROM quay.io/gravitational/debian-venti:go${GOLANG_VER}

ARG GOLANGCI_VER
ARG UID
ARG GID

# Run a script to install golangci-lint
RUN curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $GOPATH/bin ${GOLANGCI_VER}

RUN chown -R $UID:$GID /gopath
