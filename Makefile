version := $(shell git describe --tags --abbrev=0)
build := $(shell git rev-parse --short HEAD)

build:
	@go build \
		-ldflags "-X 'main.version=$(version)' -X 'main.build=$(build)'" \
		-o crmon cmd/crmon/*
