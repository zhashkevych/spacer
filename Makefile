.PHONY:
.SILENT:

build-spacer:
	go build -o ./.bin/spacer ./cmd/spacer/main.go

spacer: build-spacer
	./.bin/spacer

build-keygen:
	go build -o ./.bin/keygen ./cmd/keygen/main.go

keygen: build-keygen
	./.bin/keygen
