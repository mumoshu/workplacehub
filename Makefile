.PHONY: dep build
VERSION=$(shell git rev-parse HEAD)

dep:
		dep ensure -v

build: dep
		go build -v -ldflags "-X main.VERSION=${VERSION}"
		@mv workplacehub bin/workplacehub
