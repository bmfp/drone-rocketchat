all: build docker docker-push
deps:
	go get github.com/spf13/cobra
	go get github.com/spf13/viper
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rocket
docker:
	VERSION=`grep Version main.go | cut -d'"' -f2 > .tags`
	docker build \
		--tag registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat:$VERSION \
		--tag registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat:latest \
		--no-cache \
		--rm .
docker-push:
	docker push registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat