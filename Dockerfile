# 使用官方Go镜像作为构建环境
FROM golang:1.21.5 as builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod ./
COPY go.sum ./

# 下载依赖项
RUN go mod download

# 复制项目中的所有文件到工作目录
COPY . .

# 编译应用程序
RUN CGO_ENABLED=0 GOOS=linux go build -v -o echo-server

# 使用scratch作为最小的运行环境
FROM scratch

# 从构建器镜像中复制执行文件到当前目录
COPY --from=builder /app/echo-server .

# 设置运行时的环境变量
ENV PORT=8080

# 暴露端口
EXPOSE 8080

# 运行应用程序
CMD ["./echo-server"]