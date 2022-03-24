FROM golang:alpine as builder

COPY . /go/src/github.com/hopkiw/gce-img-resource
ENV CGO_ENABLED 0
WORKDIR /go/src/github.com/hopkiw/gce-img-resource
RUN go build -o /assets/in ./cmd/in
#RUN go build -o /assets/out ./cmd/out
RUN go build -o /assets/check ./cmd/check

FROM alpine:edge AS resource
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

FROM resource
