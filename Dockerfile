# Source: https://github.com/rebuy-de/golang-template
# Version: 2.0.2

FROM golang:1.10-alpine as builder

RUN apk add --no-cache git make

# Configure Go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Install Go Tools
RUN go get -u golang.org/x/lint/golint
RUN go get -u github.com/golang/dep/cmd/dep

COPY . /go/src/github.com/rebuy-de/kubernetes-deployment
WORKDIR /go/src/github.com/rebuy-de/kubernetes-deployment
RUN CGO_ENABLED=0 make install


FROM alpine:latest
COPY --from=builder /go/bin/kubernetes-deployment /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/kubernetes-deployment"]
