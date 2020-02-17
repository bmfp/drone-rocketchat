FROM alpine:3.11 as buildenv
RUN apk add go && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rocket

FROM alpine:3.11
COPY --from=buildenv ["rocket", "/bin/"]
RUN apk upgrade && \
    apk -Uuv add ca-certificates
ENTRYPOINT [ "/bin/rocket" ]