FROM golang:latest

WORKDIR /go/src/keys

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o keys

ENTRYPOINT ["/go/src/keys/keys"]