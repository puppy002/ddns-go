.PHONY: build clean test test-race

VERSION=0.0.1
BIN=ddns-go
DIR_SRC=.
DOCKER_CMD=docker

GO_ENV=CGO_ENABLED=0
GO_FLAGS=-ldflags="-X main.version=$(VERSION) -X 'main.buildTime=`date`' -extldflags -static"
GO=$(GO_ENV) $(shell which go)
GOROOT=$(shell `which go` env GOROOT)
GOPATH=$(shell `which go` env GOPATH)

build: $(DIR_SRC)/main.go
	@$(GO) build $(GO_FLAGS) -o $(BIN) $(DIR_SRC)

build_docker_image:
	@$(DOCKER_CMD) build -f ./Dockerfile -t ddns-go:$(VERSION) .

docker_run:
	docker run -d -p 9876:9876 --name ddns-go --restart=always  ddns-go:$(VERSION)
test:
	@$(GO) test ./...

test-race:
	@$(GO) test -race ./...

# clean all build result
clean:
	@$(GO) clean ./...
	@rm -f $(BIN)
	@rm -rf ./dist/*
