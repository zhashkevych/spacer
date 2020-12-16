.PHONY:
.SILENT:

build:
	go build -o ./.bin/spacer ./cmd/spacer/main.go

build-keygen:
	go build -o ./.bin/keygen ./cmd/keygen/main.go

keygen: build-keygen
	./.bin/keygen