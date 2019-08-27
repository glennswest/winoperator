FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 as builder
RUN go get github.com/glennswest/winoperator/winoperator
WORKDIR /go/src/github.com/glennswest/winoperator/winoperator
RUN  go get -d -v
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/winoperator
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/cli

FROM scratch
VOLUME /tmp
VOLUME /data
WORKDIR /root/
COPY --from=builder /go/bin/winoperator /go/bin/winoperator
COPY --from=builder /go/bin/cli /bin/sh
COPY commit.id commit.id
EXPOSE 8080
ENTRYPOINT ["/go/bin/winoperator"]
