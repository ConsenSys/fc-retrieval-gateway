# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

VERSION?=dev
IMAGE?=fc-retrieval-gateway

release: clean build push

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	docker build -t ${IMAGE}:${VERSION} .

# push the image to an registry
push:
	cd scripts; bash push.sh ${VERSION} ${IMAGE}:${VERSION}

utest:
	go test ./...

# remove previous images and containers
clean:
	docker rm -f ${IMAGE}:${VERSION} 2> /dev/null || true
	docker rmi -f ${IMAGE}:${VERSION} || true

cleanoldfiles:
	docker rm -f fc-retrieval-gateway-builder 2> /dev/null || true
	docker rmi -f fc-retrieval-gateway-builder || true

cleanoldfile:
	docker rm -f fc-retrieval-gateway-builder 2> /dev/null || true

.PHONY: release clean build push cleanoldfiles utest
