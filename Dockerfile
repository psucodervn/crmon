ARG BINARY=crmon
ARG DIR=/app

FROM golang:1.12.4-alpine AS builder
ARG BINARY
ARG DIR
ARG VERSION
ARG BUILD

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR $DIR

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X 'main.version=$VERSION' -X 'main.build=$BUILD'" -o $BINARY cmd/crmon/*

FROM scratch
LABEL authors="hunglm@vzota.com.vn"
ARG BINARY
ARG DIR

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder $DIR/$BINARY ./app

ENTRYPOINT ["./app"]
