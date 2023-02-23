# Go parameters
BINARY_NAME=rocket
GOCMD=go
GOBUILD=$(GOCMD) build -o $(BINARY_NAME) -v
GOBUILD_STATIC_ENVVARS=CGO_ENABLED=0 GOARCH=amd64
GOBUILD_STATIC_OPTIONS=-a -ldflags '-extldflags "-static"'
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

all: build docker docker-push
all-static: get test build-static
deps:
	go get github.com/spf13/cobra
	go get github.com/spf13/viper
build:
	$(GOBUILD)
	strip $(BINARY_NAME)
build-static:
	$(GOBUILD_STATIC_ENVVARS) $(GOBUILD) $(GOBUILD_STATIC_OPTIONS)
	strip $(BINARY_NAME)
get:
	$(GOGET)
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
docker:
	VERSION=`grep Version main.go | cut -d'"' -f2 > .tags`
	docker build \
		--tag registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat:$VERSION \
		--tag registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat:latest \
		--no-cache \
		--rm .
docker-push:
	docker push registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat