#!/bin/bash

# legacy build script to temporary fix the PR builds

cd $( dirname $0)/..

docker build .
