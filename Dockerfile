FROM golang:1.12 as builder
ENV CANNED_DIR="github.com/canned"
WORKDIR $GOPATH/src/$CANNED_DIR
COPY . .
RUN make docker_install

FROM scratch
COPY --from=builder /go/bin/canned /
ENTRYPOINT [ "/canned" ]
