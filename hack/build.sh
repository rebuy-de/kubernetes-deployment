#!/bin/bash

set -ex

cd $( dirname $0 )/..

mkdir -p target

VERSION=$(git describe --always --dirty | tr '-' '.' )

go test -v ./...

go build \
	-o target/kubernetes-deployment \
	-ldflags "-X main.version=${VERSION}" 
