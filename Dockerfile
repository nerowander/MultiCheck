# 使用 Golang 官方镜像作为构建环境
FROM golang:1.23-alpine

# 设置工作目录
WORKDIR /app

# 复制 Go 代码
COPY . .

# 编译 Go 应用
RUN apk add --no-cache git && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -o multicheck .

#EXPOSE 8080

#CMD ["/app/muticheck"]
