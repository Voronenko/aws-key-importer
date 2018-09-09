Simple console utility to import your own id_rsa.pub key (or any other) to chosen aws region, assuming you have right to do so.

# Dependencies
Uses dep utility to manage dependencies, rather than glide.
If you don't have one

`make install-dep-tool`

to install dependencies

`make deps`  or  `cd src/aws-key-importer && dep ensure`

to build binary

`make build` or `cd src/aws-key-importer && go build -o ../../bin/aws-key-importer`

#Usage 


`./aws-key-importer import notebook ~/.ssh/id_rsa.pub us-east-1`

missed values will be asked interactively 


# Roadmap:

  [ ] Key deletion
  [ ] Better error handling
  [ ] Support multiple regions in one run
