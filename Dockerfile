FROM alpine:3.11 as buildenv
RUN apk add go git && \
    cd /tmp/ && \
    git clone https://github.com/bmfp/drone-rocketchat.git && \
    cd drone-rocketchat && \
    ls -lrt . && \
    go get github.com/spf13/cobra github.com/spf13/viper && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rocket

FROM alpine:3.11
COPY --from=buildenv ["rocket", "/bin/"]
RUN apk upgrade && \
    apk -Uuv add ca-certificates
ENTRYPOINT [ "/bin/rocket" ]