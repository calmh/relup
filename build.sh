#!/bin/bash
set -euo pipefail

pkg="relup"
version="$(git describe --always)"

build() {
    name="$pkg-$version-$GOOS-$GOARCH"
    rm -rf "$name"

    gb build
    mkdir "$name"
    mv "bin/$pkg-$GOOS-$GOARCH" "$name/$pkg"
    cp LICENSE README.md "$name"
    tar zcvf "$name.tar.gz" "$name"
    rm -rf "$name"
}

GOOS=darwin  GOARCH=amd64  build
GOOS=linux   GOARCH=amd64  build
GOOS=solaris GOARCH=amd64  build
