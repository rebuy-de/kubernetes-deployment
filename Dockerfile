# Source: https://github.com/rebuy-de/golang-template
# Version: 1.2.0

FROM golang:1.8-alpine

RUN apk add --no-cache git make

# Configure Go
ENV GOPATH /go
ENV PATH /go/bin:$PATH
RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Install Glide
RUN go get -u github.com/Masterminds/glide/...

WORKDIR /go/src/github.com/Masterminds/glide

RUN git checkout v0.12.3
RUN go install

COPY . /go/src/github.com/rebuy-de/kubernetes-deployment
WORKDIR /go/src/github.com/rebuy-de/kubernetes-deployment
RUN CGO_ENABLED=0 make install

ENTRYPOINT ["/go/bin/kubernetes-deployment"]
