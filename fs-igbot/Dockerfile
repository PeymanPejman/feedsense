# TODO : Set port as ENV variable  

FROM golang:1.8

RUN apt-get update && apt-get -y install unzip && apt-get clean

EXPOSE 30200

# Install protobuf
ENV PB_VER 3.1.0

ENV PB_URL https://github.com/google/protobuf/releases/download/v${PB_VER}/protoc-${PB_VER}-linux-x86_64.zip

RUN mkdir -p /tmp/protoc && \
    curl -L ${PB_URL} > /tmp/protoc/protoc.zip && \
    cd /tmp/protoc && \
    unzip protoc.zip && \
    cp /tmp/protoc/bin/protoc /usr/local/bin && \
    cp -R /tmp/protoc/include/* /usr/local/include && \
    chmod go+rx /usr/local/bin/protoc && \
    cd /tmp && \
    rm -r /tmp/protoc

# Get the source from GitHub
RUN go get google.golang.org/grpc

# Install protoc-gen-go
RUN go get github.com/golang/protobuf/protoc-gen-go

# Make project directory
RUN mkdir -p /go/src/feedsense/fs-igbot

# Set working directory
WORKDIR /go/src/feedsense/fs-igbot/src

# Add source code directory
ADD . / /go/src/feedsense/fs-igbot/

# Build go binaries
RUN go install -v feedsense/fs-igbot/src

# Run go binary
CMD ["src"]