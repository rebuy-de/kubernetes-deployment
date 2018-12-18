# Source: https://github.com/rebuy-de/golang-template
# Version: 2.0.2

FROM golang:1.11-alpine as builder

RUN apk add --no-cache git make curl openssl 

# Configure Go
ENV GOPATH=/go PATH=/go/bin:$PATH CGO_ENABLED=0 GO111MODULE=on
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Install Go Tools
RUN GO111MODULE= go get -u golang.org/x/lint/golint

# Install Linkerd
RUN curl -Lo /usr/local/bin/linkerd https://github.com/linkerd/linkerd2/releases/download/stable-2.1.0/linkerd2-cli-stable-2.1.0-linux
RUN chmod +x /usr/local/bin/linkerd
RUN linkerd version --client --api-addr="localhost"

COPY . /src
WORKDIR /src
RUN make build

FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /usr/local/bin/linkerd /usr/local/bin/
COPY --from=builder /src/dist/* /usr/local/bin/
