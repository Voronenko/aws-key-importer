SHELL := /bin/bash
PATH := $(CURDIR)/.gopath/bin:$(PATH)
VERSION:=$(shell cat version.txt)
PACKAGE  = github.com/voronenko/aws_key_importer
GOPATH   = $(CURDIR)/.gopath
BASE     = $(GOPATH)/src/$(PACKAGE)

$(BASE):
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

local_bin_deps: $(BASE)
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@go get -u github.com/golang/dep/cmd/dep

bin_deps: $(BASE)
	@curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(GOPATH)/bin

install-githubrelease-tool: $(BASE)
	go get github.com/aktau/github-release

deps: $(BASE)
	cd $(BASE) && \
	dep ensure

lint: $(BASE)
	cd $(BASE)/cmd/aws-key-importer && golangci-lint run

build: deps
	cd $(BASE)/cmd/aws-key-importer && go build -o ../../out/aws-key-importer

upload:
	github-release upload  --user voronenko --repo aws-key-importer --tag $(VERSION)  --name "aws-key-importer-linux-amd64" --file out/aws-key-importer
