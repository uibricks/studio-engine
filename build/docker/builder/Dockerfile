FROM golang:1.17-alpine AS builder

# Install `git` binary
RUN apk update
RUN apk add git build-base

## Install go tools
RUN go install github.com/go-task/task/v3/cmd/task@latest
WORKDIR /service

COPY go.mod go.sum ./
RUN go mod download