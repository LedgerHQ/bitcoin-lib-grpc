FROM golang:1.15-alpine as builder

# System setup
RUN apk update && apk add git curl build-base autoconf automake libtool

# Install protoc
ENV PROTOBUF_VERSION 3.13.0
ENV PROTOBUF_URL https://github.com/google/protobuf/releases/download/v$PROTOBUF_VERSION/protobuf-cpp-$PROTOBUF_VERSION.tar.gz
RUN curl -L -o /tmp/protobuf.tar.gz $PROTOBUF_URL
WORKDIR /tmp/
RUN tar xvzf protobuf.tar.gz
WORKDIR /tmp/protobuf-$PROTOBUF_VERSION
RUN mkdir /tmp/protobuf
RUN ./autogen.sh && ./configure --prefix=/tmp/protobuf && make && make install
RUN ln -s /tmp/protobuf/bin/protoc /usr/local/bin/protoc

# Install mage
RUN mkdir /tmp/mage
WORKDIR /tmp/mage
RUN git clone https://github.com/magefile/mage && cd mage && go run bootstrap.go

# Install gRPC and protobuf tools
RUN go get google.golang.org/grpc
RUN go get github.com/golang/protobuf/protoc-gen-go

# Set the current working directory inside the container
WORKDIR /app

COPY . .

# Build the Go app
RUN mage -v build

# Start fresh from a smaller image
FROM alpine

COPY --from=builder /app/lbs /app/lbs

ENV GRPC_GO_LOG_SEVERITY_LEVEL info
ENV GRPC_GO_LOG_VERBOSITY_LEVEL 1

EXPOSE 50051

CMD ["/app/lbs"]
