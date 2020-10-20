.PHONY: run build

upx := $(shell which upx)

run:
	go run .

build:
	go build .
ifdef upx
	upx $$(basename $$(pwd))
endif
