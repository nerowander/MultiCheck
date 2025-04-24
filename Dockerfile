FROM golang:1.23-alpine

WORKDIR /app

COPY . .

RUN apk add --no-cache git && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -o multicheck .
