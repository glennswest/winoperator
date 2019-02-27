FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 
RUN git clone https://github.com/glennswest/winoperator /go/src/github.com/glennswest/winoperator
WORKDIR /go/src/github.com/glennswest/winoperator/winoperator
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 

FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
COPY --from=0 /go/src/github.com/glennswest/winoperator/winoperator /usr/bin/winoperator
ENTRYPOINT ["/usr/bin/winoperator"]
