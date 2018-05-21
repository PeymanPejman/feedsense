FROM golang:1.8

RUN mkdir -p /go/src/feedsense

WORKDIR /go/src/feedsense/src

ADD . / /go/src/feedsense/

RUN go install -v feedsense/src

CMD ["src"]

# Set port as ENV variable here 