#!/bin/bash

REPO=github.com/rebuy-de/kubernetes-deployment
ROOT=$(readlink -f $( dirname $0)/..)
COMMAND=${1:-all}

set -ex

docker run \
	--rm \
	-v "${ROOT}:/go/src/${REPO}" \
	074509403805.dkr.ecr.eu-west-1.amazonaws.com/rebuy-base-image-golang:latest \
    /tools/build.sh ${COMMAND} ${REPO}
