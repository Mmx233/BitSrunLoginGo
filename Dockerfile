FROM golang:alpine as builder

RUN go env -w CGO_ENABLED=0

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags '-extldflags "-static -fpic" -s -w' -o runner ./cmd/bitsrun

FROM alpine:latest

RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache ca-certificates &&\
    rm -rf /var/cache/apk/*

COPY --from=builder /build/runner /usr/bin/bitsrun
WORKDIR /data

ENTRYPOINT [ "/usr/bin/bitsrun" ]
