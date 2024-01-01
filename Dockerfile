FROM golang:1.20.0-alpine3.17 as BUILDER
ADD . /go/src/github.com/dedup
WORKDIR /go/src/github.com/dedup
RUN cd /go/src/github.com/dedup && \
    go build -o dedup

FROM alpine:3.17.1
WORKDIR /etc/dedup
COPY --from=BUILDER /go/src/github.com/dedup/dedup /usr/local/bin

CMD ["/usr/local/bin/dedup"]