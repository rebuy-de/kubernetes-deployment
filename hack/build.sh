#!/bin/bash

set -ex

cd $( dirname $0 )/..

mkdir -p target

VERSION=$(git describe --always --dirty | tr '-' '.' )

glide install

go test -v $(glide nv)

go build \
	-o target/kubernetes-deployment \
	-ldflags "-X main.version=${VERSION}" 
