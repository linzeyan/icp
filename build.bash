#!/bin/bash

set -ex

VERSION="$(git describe --tags --exact-match)"
PACKAGE="icp"
TARGET="bin"

build() {
    GOOS=darwin GOARCH=amd64
    go build -a -trimpath -o ${TARGET}/${PACKAGE}_${VERSION}_${GOOS}_${GOARCH} icp.go
    # Linux
    GOOS=linux GOARCH=amd64
    go build -a -trimpath -o ${TARGET}/${PACKAGE}_${VERSION}_${GOOS}_${GOARCH} icp.go
    GOOS=linux GOARCH=arm64
    go build -a -trimpath -o ${TARGET}/${PACKAGE}_${VERSION}_${GOOS}_${GOARCH} icp.go
    # Windows
    GOOS=windows GOARCH=amd64
    go build -a -trimpath -o ${TARGET}/${PACKAGE}_${VERSION}_${GOOS}_${GOARCH} icp.go
}

clean() {
    go clean
    rm -rf ${TARGET}
}

$1
