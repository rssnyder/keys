FROM golang:latest

ARG TARGETOS
ARG TARGETARCH

WORKDIR /go/src/keys

COPY . .

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o keys

ENTRYPOINT ["/go/src/keys/keys"]