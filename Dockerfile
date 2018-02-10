FROM golang:1.8
LABEL maintainer "Albert Shin <shina2@rpi.edu>"

ENV PATH ${GOPATH}/bin:$PATH

WORKDIR ${GOPATH}/src/github.com/albshin/tutescrew

COPY . .

RUN go build -o tutescrew cmd/tutescrew/main.go

EXPOSE 8080

CMD ["./tutescrew"]