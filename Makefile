.PHONY: default build clean build-image push-image

BINARY = demo-controller

DOCKER_REPO = YOUR_URL

VER := $(shell grep "const version = " main.go | cut -f2 -d'"')
GIT := $(shell git rev-parse --short HEAD)
DOCKER_TAG := $(VER)-$(GIT)

GOCMD = go
GOFLAGS ?= $(GOFLAGS:)
LDFLAGS =

default: build

build:
	"$(GOCMD)" build ${GOFLAGS} ${LDFLAGS} -o "${BINARY}"

build-image:
	@docker build -t ${DOCKER_REPO}:${DOCKER_TAG} .

push-image: build-image
	@docker push ${DOCKER_REPO}:${DOCKER_TAG}

clean:
	"$(GOCMD)" clean -i
