# 使用官方的 Go 镜像作为基础镜像
FROM golang:1.23-alpine AS build-env

# 设置工作目录
WORKDIR /app

# 将当前目录下的所有文件复制到容器的 /app 目录下
COPY . .

ENV GOPROXY='https://goproxy.cn,direct'
# 构建 Go 应用
RUN cd server && go build -o main .

# 使用更小的基础镜像来减少最终镜像的大小
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的二进制文件到新的镜像中
COPY --from=build-env /app/server/main .

# 暴露服务端口
EXPOSE 50051

# 设置容器启动时运行的命令
CMD ["./main"]
