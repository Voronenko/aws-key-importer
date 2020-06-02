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

deps:
	go mod download
  go mod vendor

lint:
	cd cmd/aws-key-importer && golangci-lint run

deps:
	go mod download
	go mod vendor

build: deps
	cd $(BASE)/cmd/aws-key-importer && go build -o ../../dist/aws-key-importer
