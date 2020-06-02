VERSION:=$(shell cat version.txt)

install-golint-tool:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(GOPATH)/bin

lint:
	cd cmd/aws-key-importer && golangci-lint run

deps:
	go mod download
	go mod vendor

build: deps
	cd cmd/aws-key-importer && go build -o ../../dist/aws-key-importer
