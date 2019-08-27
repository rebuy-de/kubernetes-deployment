# Source: https://github.com/rebuy-de/golang-template

FROM golang:1.12-alpine as builder

RUN apk add --no-cache git make curl openssl

# Configure Go
ENV GOPATH=/go PATH=/go/bin:$PATH CGO_ENABLED=0 GO111MODULE=on
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Install Go Tools
RUN GO111MODULE= go get -u golang.org/x/lint/golint

# Install Linkerd
RUN set -x \
 && curl -Lo /usr/local/bin/linkerd https://github.com/linkerd/linkerd2/releases/download/stable-2.5.0/linkerd2-cli-stable-2.5.0-linux \
 && sha256sum /usr/local/bin/linkerd \
 && echo "2a4e0d3982c67a2cdd39be8c1a4de1ac2f70cba4b6888a8e57553d900432b599  /usr/local/bin/linkerd" | sha256sum -c \
 && chmod +x /usr/local/bin/linkerd \
 && linkerd version --client --api-addr="localhost"

# Install kubectl
RUN set -x \
 && curl -O https://storage.googleapis.com/kubernetes-release/release/v1.15.3/bin/linux/amd64/kubectl \
 && mv kubectl /usr/local/bin/kubectl \
 && chmod 755 /usr/local/bin/kubectl \
 && kubectl version --client

COPY . /src
WORKDIR /src
RUN set -x \
 && make test \
 && make build \
 && cp --dereference /src/dist/* /usr/local/bin/

RUN set -x \
 && kubernetes-deployment version

FROM alpine:latest
RUN apk add --no-cache ca-certificates

COPY --from=builder /usr/local/bin/* /usr/local/bin/

RUN adduser -D k26r
USER k26r
