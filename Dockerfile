FROM golang:1.10.2 as builder
WORKDIR /go/src/build
COPY vendor/ ./vendor/
COPY pkg/ ./
COPY *.go ./
RUN go test
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o demo-controller .

FROM alpine:3.7  
WORKDIR /
COPY --from=builder /go/src/build/demo-controller .
CMD ["/demo-controller"] 