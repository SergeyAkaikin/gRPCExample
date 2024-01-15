FROM golang:1.21-alpine as build-env

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

RUN apk --no-cache update && apk add --no-cache git make zip && \
    mkdir /opt/protoc && cd /opt/protoc && wget https://github.com/protocolbuffers/protobuf/releases/download/v25.2/protoc-25.2-linux-x86_64.zip && \
    unzip protoc-25.2-linux-x86_64.zip

ENV PATH=$PATH:$GOPATH/bin:/opt/protoc/bin

RUN mkdir -p /go/src/test
COPY ProductInfo.proto /go/src/test

WORKDIR /go/src/test

COPY go.mod .
COPY go.sum .

RUN go mod download && mkdir -p ./gen/go

RUN protoc -I . ./ProductInfo.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go --go-grpc_opt=paths=source_relative
RUN mkdir -p grpc-server

COPY ./service/main.go .

RUN go build -tags musl -o bin/test ./main.go

FROM alpine:3.18.5

COPY --from=build-env /go/src/test/bin/test /usr/local/bin/test

EXPOSE 50051
CMD ["test"]