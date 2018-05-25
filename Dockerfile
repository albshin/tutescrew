
FROM golang:1.10 AS builder
LABEL maintainer "Albert Shin <shina2@rpi.edu>"

# Download and install the latest release of dep
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep

# Copy the code from the host and compile it
WORKDIR $GOPATH/src/github.com/albshin/tutescrew
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /app .

FROM alpine:latest
COPY --from=builder /app ./

RUN apk --no-cache add ca-certificates

ENV TUTESCREW_TOKEN="MzQ5NjMyMTkwMzk4ODU3MjE3.Den-7Q.fzF_3v-5ruO4Gl0dGKq1SUf37CM"
ENV TUTESCREW_CAS_AUTHURL="https://cas-auth.rpi.edu/cas/login"
ENV TUTESCREW_CAS_REDIRECTURL="http://localhost:8080/auth/cas"

ENTRYPOINT ["./app"]

EXPOSE 8080