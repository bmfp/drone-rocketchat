all: build docker docker-push
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rocket
docker:
	VERSION=`./rocket --version | cut -d" " -f3`
	docker build \
		--tag registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat:$VERSION \
		--tag registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat:latest \
		--no-cache \
		--rm .
docker-push:
	docker push registry.ver.bmfp.fr/drone-rocketchat/drone-rocketchat