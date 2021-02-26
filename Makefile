# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

VERSION?=dev
IMAGE?=consensys/fc-retrieval-gateway

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

# User `make dev arg=--build` to rebuild
dev:
	make -s env-soft
	docker-compose -f docker-compose.dev.yml up $(arg)

build-dev:
	go build -v cmd/gateway/main.go

env-soft:
	make env-compare
	if [ -s .env ]; then \
		echo ".env already exists"; \
	else \
		echo ".env does not exists, create a .env file"; \
		cp .env.example .env; \
	fi

env-compare:
	@cmp -s .env .env.example; \
	RETVAL=$$?; \
	if [ $$RETVAL -eq 0 ]; then \
		echo ".env and .env.example are the same"; \
	else \
		echo ".env and .env.example are NOT the same"; \
	fi