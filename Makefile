CLIENT_IMAGE_VERSION ?= 0.0.1
SERVER_IMAGE_VERSION ?= 0.0.1

_GO_BUILD = \
	docker run --rm \
		-v $(shell pwd):/home/engine/repo \
		--workdir /home/engine/repo \
		golang:1.16-alpine \
		go build -v -o

build/client:
	$(_GO_BUILD) ./out/bin/wow-client ./cmd/wow-client/main.go

build/server:
	$(_GO_BUILD) ./out/bin/wow-server ./cmd/wow-server/main.go


docker/client-image: build/client
	docker build \
		-f docker/client/Dockerfile \
		-t go-wow-client:${CLIENT_IMAGE_VERSION} \
		.

docker/server-image: build/server
	docker build \
		-f docker/server/Dockerfile \
		-t go-wow-server:${SERVER_IMAGE_VERSION} \
		.

docker/images: docker/client-image docker/server-image

docker/network:
	docker network create go-wow-network || exit 0

start/server: docker/network
	docker run \
		--rm -d \
		--cpus 1 \
		--oom-kill-disable \
		--memory 300M \
		--name go-wow-server \
		--network go-wow-network \
		go-wow-server:${SERVER_IMAGE_VERSION}

stop/server:
	docker stop go-wow-server

start/client:
		# echo wow-client -addr go-wow-server:1024 -count 25000
		docker run \
		-it --rm \
		--network go-wow-network \
		go-wow-client:${CLIENT_IMAGE_VERSION} \
		sh

start/client-spam:
		docker run -d \
		-it --rm \
		--network go-wow-network \
		go-wow-client:${CLIENT_IMAGE_VERSION} \
		wow-client -addr go-wow-server:1024 -count 25000
