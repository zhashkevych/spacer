.PHONY:
.SILENT:

build:
	go build -o ./.bin/spacer

run: build
	./.bin/spacer -bucket jewerly -endpoint ams3.digitaloceanspaces.com -db_port 5436