#!/bin/bash
set -euo pipefail

pkg="relup"
version="$(git describe --always)"

build() {
    name="$pkg-$version-$GOOS-$GOARCH"
    rm -rf "$name"

    go build
    mkdir "$name"
    mv "$pkg" "$name/$pkg"
    cp LICENSE README.md "$name"
    tar zcvf "$name.tar.gz" "$name"
    rm -rf "$name"
}

GOOS=darwin  GOARCH=amd64  build
GOOS=linux   GOARCH=amd64  build
GOOS=solaris GOARCH=amd64  build
