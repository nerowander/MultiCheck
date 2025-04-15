# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.23-alpine

WORKDIR /app

COPY ./PocScan .
COPY ./WebScan .
COPY ./common .
COPY ./config .
COPY ./go.mod .

RUN go mod tidy && \
    go build -o pocscan .

EXPOSE 8080

CMD ["./pocscan"]