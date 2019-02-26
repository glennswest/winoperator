FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 AS builder
WORKDIR /go/src/github.com/glennswest/winoperator
COPY . .
RUN WHAT=winoperator ./hack/build-go.sh

FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
COPY --from=builder /go/src/github.com/glennswest/winoperator/_output/linux/amd64/winoperator /usr/bin/
ENTRYPOINT ["/usr/bin/winoperator"]
