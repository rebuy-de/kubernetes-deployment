#!/bin/bash

set -ex
REPO=github.com/rebuy-de/kubernetes-deployment
ROOT=$(dirname $(dirname $0))
COMMAND=${1:-all}

docker run \
	--rm \
	-v "${PWD}:/go/src/${REPO}" \
	074509403805.dkr.ecr.eu-west-1.amazonaws.com/rebuy-base-image-golang:latest \
    /tools/build.sh ${COMMAND} ${REPO}
