.PHONY:
.SILENT:

build:
	go build -o ./.bin/spacer ./cmd/main.go

spacer: build
	./.bin/spacer