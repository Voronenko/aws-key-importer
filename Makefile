install-dep-tool:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

install-golint-tool:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b $(GOPATH)/bin

install-githubrelease-tool:
	go get github.com/aktau/github-release

deps:
	cd src/aws-key-importer && dep ensure

lint:
	cd src/aws-key-importer && golangci-lint run

build: deps
	cd src/aws-key-importer && go build -o ../../bin/aws-key-importer

upload:
	github-release upload  --user voronenko --repo aws-key-importer --tag 0.1.0  --name "aws-key-importer-linux-amd64" --file bin/aws-key-importer
