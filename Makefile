# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

VERSION?=v1
IMAGE?=fc-retrieval-gateway:dev


release: clean build push

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	docker build -t ${IMAGE} .

# push the image to an registry
push:
	cd scripts; bash push.sh ${VERSION} ${IMAGE}

utest:
	go test ./...

# remove previous images and containers
clean:
	docker rm -f ${IMAGE} 2> /dev/null || true
	docker rmi -f ${IMAGE} || true

cleanoldfiles:
	docker rm -f fc-retrieval-gateway-builder 2> /dev/null || true
	docker rmi -f fc-retrieval-gateway-builder || true

cleanoldfile:
	docker rm -f ${REGISTRY}fc-retrieval-gateway-builder 2> /dev/null || true

.PHONY: release clean build push cleanoldfiles utest
