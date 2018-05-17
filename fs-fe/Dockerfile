FROM golang:1.8

RUN mkdir -p /feedsense

WORKDIR /feedsense

ADD . /feedsense

RUN go build ./hello.go

CMD ["./hello"]
