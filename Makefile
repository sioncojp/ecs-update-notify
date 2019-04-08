REVISION := $(shell git describe --always)
DATE := $(shell date +%Y-%m-%dT%H:%M:%S%z)
LDFLAGS	:= -ldflags="-X \"main.Revision=$(REVISION)\" -X \"main.BuildDate=${DATE}\""

.PHONY: build-cross dist build clean run help

name		:= ecs-update-notify
linux_name	:= $(name)-linux-amd64
darwin_name	:= $(name)-darwin-amd64

dist: build-docker ## create .tar.gz linux & darwin to /bin
	cd bin && tar zcvf $(linux_name).tar.gz $(linux_name) && rm -f $(linux_name)
	cd bin && tar zcvf $(darwin_name).tar.gz $(darwin_name) && rm -f $(darwin_name)

build: ## go build
	go build -o bin/$(name) $(LDFLAGS) cmd/$(name)/*.go

build-cross: ## create to build for linux & darwin to bin/
	GOOS=linux GOARCH=amd64 go build -o bin/$(linux_name) $(LDFLAGS) cmd/$(name)/*.go
	GOOS=darwin GOARCH=amd64 go build -o bin/$(darwin_name) $(LDFLAGS) cmd/$(name)/*.go

build-docker: ## go build on Docker
	docker run --rm -v "$(PWD)":/go/src/github.com/sioncojp/$(name) -w /go/src/github.com/sioncojp/$(name) golang:latest bash build.sh

test: ## go test
	go test -v $$(go list ./... | grep -v /vendor/)

clean: ## remove bin/*
	rm -f bin/*

run: ## go run
	go run cmd/$(name)/main.go -c examples/config.toml

lint: ## go lint ignore vendor
	golint $(go list ./... | grep -v /vendor/)

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-22s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)
