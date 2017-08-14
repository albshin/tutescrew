FROM golang:1.8.3-alpine3.6

LABEL Name=tutescrew Version=0.0.1

RUN mkdir -p /go/src/github.com/albshin/tutescrew

ADD . /go/src/github.com/albshin/tutescrew

RUN go install github.com/albshin/tutescrew

ENTRYPOINT /go/bin/tutescrew

EXPOSE 8080