# Copyright 2020 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.14 as builder

WORKDIR /go/src/custom-plugin

COPY ./server/main.go .
COPY ./server/go.mod .
COPY ./server/go.sum . 


ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -a -ldflags '-s -w -extldflags "-static"' -o /go/bin/custom-plugin

FROM envoyproxy/envoy-alpine:v1.15.0

COPY envoy.yaml /etc/envoy/envoy.yaml
COPY --from=builder /go/bin/custom-plugin .
COPY startup.sh startup.sh

RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

EXPOSE 8000
EXPOSE 8080

CMD ./startup.sh