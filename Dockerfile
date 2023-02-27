FROM alpine:3.17
RUN apk add --update-cache --no-cache --upgrade --virtual buildenv go git && \
    apk add --update-cache --no-cache --upgrade ca-certificates && \
    cd /tmp/ && \
    git clone https://github.com/bmfp/drone-rocketchat.git && \
    cd drone-rocketchat && \
    go get -v github.com/spf13/cobra github.com/spf13/viper && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /bin/rocket && \
    strip /bin/rocket && \
    apk del buildenv && \
    rm -rf /tmp/drone-rocketchat /root/go /root/.cache 

ENTRYPOINT [ "/bin/rocket" ]
