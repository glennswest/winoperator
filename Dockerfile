FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 as builder
RUN go get github.com/glennswest/winoperator/winoperator
WORKDIR /go/src/github.com/glennswest/winoperator/winoperator
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' .;ls -l

FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/glennswest/winoperator/winoperator/winoperator /root/winoperator
RUN MKDIR /tmp
ENTRYPOINT ["/root/winoperator"]
