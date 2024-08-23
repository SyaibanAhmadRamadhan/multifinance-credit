#####################
ARG GO_VERSION=1.22

FROM golang:${GO_VERSION} AS builder

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Install dependencies
RUN apt-get update && \
   apt-get install -y \
   nodejs \
   npm \
   make \
   gcc \
   git \
   unzip \
   wget

RUN mkdir -p /api

WORKDIR /api

# Install dependencies
COPY . .
RUN make install

# Install go mod dependencies

RUN go mod download

RUN make generate

RUN echo "Building API" && \
   go build -o main ./main.go

#
FROM alpine:latest

RUN mkdir -p /api

# Install glibc

#using gcompat or force install glibc(nimpa musl libc).
#in case using gcompat its easy installation
RUN apk update && apk add --no-cache gcompat libstdc++

COPY --from=builder /api/main /api
COPY --from=builder /api/env.json /api

WORKDIR /api

# Ensure the main file has execution permissions
RUN chmod +x main

ENV PORT=80

EXPOSE 80

ENTRYPOINT ["./main"]