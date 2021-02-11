#!/bin/bash
echo "****************************************"
echo "*** Push docker image to Docker Hub  ***"
echo "****************************************"

# Docker image name on docker hub
IMAGE_NAME=consensys/fc-retrieval-gateway

VERSION=$1
echo "Gateway version: $VERSION"
IMAGE=$2
echo "Gateway image: $IMAGE"


GATEWAY_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "Gateway repo branch: $GATEWAY_BRANCH"

TAG="develop-$GATEWAY_BRANCH"
echo "TAG: $TAG"
docker tag $IMAGE $IMAGE_NAME:$TAG
