.PHONY: run build

upx := $(shell which upx)
name := "aws-helper"

run:
	go run .

build:
	go build \
		-o ${name} \
		-ldflags "-X main.version=$$(cat version.json | jq -r .version)" \
		.
ifdef upx
	upx $$(basename $$(pwd))
endif
