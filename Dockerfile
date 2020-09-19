FROM golang:1.14 as builder
ADD . /go/src/custom-plugin 
WORKDIR /go/src/custom-plugin
COPY . /go/src/custom-plugin
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -a -ldflags '-s -w -extldflags "-static"' -o /go/bin/custom-plugin

FROM envoyproxy/envoy-alpine:v1.14.1

COPY envoy-local.yaml /etc/envoy/envoy.yaml
COPY --from=builder /go/bin/custom-plugin .
COPY startup.sh startup.sh

EXPOSE 9000
EXPOSE 9080

CMD ./startup.sh