REVISION := $(shell git describe --always)
LDFLAGS	 := -ldflags="-X \"main.Revision=$(REVISION)\""

.PHONY: build-cross dist build clean run help go/* build

name		    := ecs-update-notify
linux_name  	:= $(name)-linux-amd64
darwin_name   	:= $(name)-darwin-amd64
darwin_arm_name	:= $(name)-darwin-arm64

go_version := $(shell cat $(realpath .go-version))
bindir     := $(realpath bin/)
go         := $(bindir)/go/bin/go
arch       := $(shell arch)

help:
	@awk -F ':|##' '/^[^\t].+?:.*?##/ { printf "\033[36m%-22s\033[0m %s\n", $$1, $$NF }' $(MAKEFILE_LIST)

### go install
go/install: file         = go.tar.gz
go/install: download_url = https://golang.org/dl/go$(go_version).darwin-$(arch).tar.gz
go/install:
# If you have a different version, delete it.
	@if [ -f $(go) ]; then \
		$(go) version | grep -q "$(go_version)" || rm -f $(bindir)/go; \
	fi

# If the file is not there, download it.
	@if [ ! -f $(go) ]; then \
		curl -L -fsS --retry 2 -o $(file) $(download_url) && \
		tar zxvf $(file) -C $(bindir) && rm -f $(file); \
	fi

dist: build-docker ## create .tar.gz linux & darwin to /bin
	cd bin && tar zcvf $(linux_name).tar.gz $(linux_name) && rm -f $(linux_name)
	cd bin && tar zcvf $(darwin_name).tar.gz $(darwin_name) && rm -f $(darwin_name)
	cd bin && tar zcvf $(darwin_arm_name).tar.gz $(darwin_arm_name) && rm -f $(darwin_arm_name)

test: ## go test
	$(go) test -v $$(go list ./... | grep -v /vendor/)

clean: ## remove bin/*
	rm -f bin/*

lint: ## go lint ignore vendor
	golint $($(go) list ./... | grep -v /vendor/)


build: go/install ## build
	$(go) build -o bin/$(name) cmd/$(name)/*.go

build/cross: go/install ## create to build for linux & darwin to bin/
	GOOS=linux GOARCH=amd64 $(go) build -o bin/$(linux_name) $(LDFLAGS) cmd/$(name)/*.go
	GOOS=darwin GOARCH=amd64 $(go) build -o bin/$(darwin_name) $(LDFLAGS) cmd/$(name)/*.go
	GOOS=darwin GOARCH=arm64 $(go) build -o bin/$(darwin_arm_name) $(LDFLAGS) cmd/$(name)/*.go

build-docker: ## go build on Docker
	docker run --rm -v "$(PWD)":/go/src/github.com/sioncojp/$(name) -w /go/src/github.com/sioncojp/$(name) golang:$(go_version) bash build.sh

run: ## go run
	$(go) run cmd/$(name)/main.go -c examples/config.toml
