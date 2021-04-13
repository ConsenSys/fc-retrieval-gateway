# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

VERSION?=dev
IMAGE?=consensys/fc-retrieval-gateway
COV?=80

release: clean build tag

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	docker build -t ${IMAGE}:${VERSION} .

# push the image to an registry
push:
	cd scripts; bash push.sh ${VERSION} ${IMAGE}:${VERSION}

tag:
	cd scripts; bash tag.sh ${VERSION} ${IMAGE}:${VERSION}

useremote:
	cd scripts; bash use-remote-repos.sh


utest:
	go test ./...

coverage:
	bash ./scripts/coverage.sh $(COV)

# remove previous images and containers
clean:
	docker rm -f ${IMAGE}:${VERSION} 2> /dev/null || true
	docker rmi -f ${IMAGE}:${VERSION} || true

cleanoldfiles:
	docker rm -f fc-retrieval-gateway-builder 2> /dev/null || true
	docker rmi -f fc-retrieval-gateway-builder || true

cleanoldfile:
	docker rm -f fc-retrieval-gateway-builder 2> /dev/null || true

.PHONY: release clean build push cleanoldfiles utest coverage

# User `make dev arg=--build` to rebuild
dev:
	./scripts/make-env-file/make-env-file -source=.env.example -dest=.env
	docker-compose -f docker-compose.dev.yml up $(arg)

build-dev:
	go build -v cmd/gateway/main.go
