# 使用官方 Golang 镜像作为基础镜像
FROM golang:1.23-alpine

WORKDIR /app

#COPY ./config .
#COPY ./InformationScan .
#COPY ../common .
#COPY ./WebScan .
#COPY ./go.mod .
COPY . .

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -o infomationscan .

EXPOSE 8080

CMD ["./infomationscan"]