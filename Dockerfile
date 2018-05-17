FROM golang:1.10.2 as builder
WORKDIR /go/src/github.com/piontec/k8s-demo-controller
COPY vendor/ ./vendor/
COPY pkg/ ./pkg/
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demo-controller .

FROM alpine:3.7  
WORKDIR /
COPY --from=builder /go/src/github.com/piontec/k8s-demo-controller/demo-controller .
CMD ["/demo-controller"] 