install-dep-tool:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

deps:
	cd src/aws-key-importer && dep ensure

build: deps
	cd src/aws-key-importer && go build -o ../../bin/aws-key-importer
