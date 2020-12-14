.PHONY:
.SILENT:

build:
	go build -o ./.bin/spacer

run: build
	./.bin/spacer