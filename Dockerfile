# References the official image for Go as the base image.
FROM golang:1.16-alpine

ADD . /go/src/dating
WORKDIR /go/src/dating
RUN go get dating
RUN go mod vendor
RUN go install
EXPOSE 8080

ENTRYPOINT ["/go/bin/dating"]