# BUILD
FROM golang:latest as builder

# copy
WORKDIR /src/
COPY . /src/

# build
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -tags netgo \
    -ldflags '-w -extldflags "-static"' \
    -mod vendor \
    -o eventstore

# CERTS
FROM alpine:latest as certs
RUN apk --update add ca-certificates

# RUN
FROM scratch
COPY --from=builder /src/eventstore .
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ENTRYPOINT ["/eventstore"]