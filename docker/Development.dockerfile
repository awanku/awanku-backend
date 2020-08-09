FROM golang:1.14.4-alpine
RUN apk add --no-cache build-base xz curl ca-certificates make postgresql-client
RUN curl -o /tmp/watchexec.tar.xz -sL https://github.com/watchexec/watchexec/releases/download/1.13.1/watchexec-1.13.1-x86_64-unknown-linux-musl.tar.xz && \
    tar -xJ -f /tmp/watchexec.tar.xz -C /tmp && \
    mv /tmp/watchexec-1.13.1-x86_64-unknown-linux-musl/watchexec /usr/bin/watchexec
