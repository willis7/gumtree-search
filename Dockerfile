FROM golang:1.9 as builder
WORKDIR /go/src/github.com/willis7/gumtree-searcher
COPY ./ ./
RUN go get golang.org/x/net/html
RUN CGO_ENABLED=0 go install -a -tags gumtree-searcher -ldflags '-extldflags "-static"'
RUN ldd /go/bin/gumtree-searcher | grep -q "not a dynamic executable"

FROM alpine
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /go/bin/gumtree-searcher /gumtree

ENTRYPOINT ["/gumtree"]
