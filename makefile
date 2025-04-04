NAME       := fuzzy-guacamole
GO_FLAGS   ?=
GO_TAGS	   ?= netgo
CGO_ENABLED?=0
OUTPUT_BIN ?= execs/${NAME}
PACKAGE    := github.com/Aearsears/$(NAME)
GIT_REV    ?= $(shell git rev-parse --short HEAD)
SOURCE_DATE_EPOCH ?= $(shell date +%s)
ifeq ($(shell uname), Darwin)
DATE       ?= $(shell TZ=UTC date -j -f "%s" ${SOURCE_DATE_EPOCH} +"%Y-%m-%dT%H:%M:%SZ")
else
DATE       ?= $(shell date -u -d @${SOURCE_DATE_EPOCH} +"%Y-%m-%dT%H:%M:%SZ")
endif
VERSION    ?= v0.0.0
IMG_NAME   := Aearsears/$(NAME)
IMAGE      := ${IMG_NAME}:${VERSION}

default: help

test:   ## Run all tests
	@go clean --testcache && go test ./...

cover:  ## Run test coverage suite
	@go test ./... --coverprofile=cov.out
	@go tool cover --html=cov.out

build:  ## Builds the CLI
	@CGO_ENABLED=${CGO_ENABLED} go build ${GO_FLAGS} \
	-ldflags "-w -s -X ${PACKAGE}/cmd.version=${VERSION} -X ${PACKAGE}/cmd.commit=${GIT_REV} -X ${PACKAGE}/cmd.date=${DATE}" \
	-a -tags=${GO_TAGS} -o ${OUTPUT_BIN} main.go

localstack:    
	@docker pull localstack/localstack
	@docker run -d -p 4566:4566 -p 4571:4571 --name localstack \
  	-e SERVICES=s3 \
  	-e DEBUG=1 \
  	-v "$HOME/localstack:/tmp/localstack" \
  	localstack/localstack		

dev:
	@go run main.go --config ./configs/dev.yaml --log-level debug --log-format text --log-output stdout

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":[^:]*?## "}; {printf "\033[38;5;69m%-30s\033[38;5;38m %s\033[0m\n", $$1, $$2}'