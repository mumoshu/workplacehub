.PHONY: dep build image container
VERSION=$(shell git rev-parse HEAD)

dep:
		dep ensure -v

build: dep
		go build -v -ldflags "-X main.VERSION=${VERSION}"
		@mv workplacehub bin/workplacehub

image:
	docker build -t workplacehub:canary .

container: image
	docker run -it -p 8010:8010 workplacehub:canary
