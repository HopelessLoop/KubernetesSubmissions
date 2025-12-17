# -------------------------------------------------------------------
# 第一阶段：构建 (Builder)
# 使用官方 Go 镜像作为构建环境
# -------------------------------------------------------------------
FROM golang:1.24-alpine AS builder

# 设置工作目录
WORKDIR /app

# 1. 先拷贝依赖文件，利用 Docker 缓存层
COPY go.mod go.sum ./

RUN go mod download

# 2. 拷贝源码并编译
COPY . .

# 关键参数：
# CGO_ENABLED=0: 禁用 CGO，确保生成纯静态二进制文件（否则在 scratch/alpine 中可能运行失败）
# -ldflags="-s -w": 去掉符号表和调试信息，减小体积
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o myapp .

# -------------------------------------------------------------------
# 第二阶段：运行 (Runner)
# 选择一个极其精简的镜像
# -------------------------------------------------------------------
FROM alpine:latest

WORKDIR /usr/src/app

# 从第一阶段拷贝编译好的二进制文件
COPY --from=builder /app/myapp .

ENV PORT=8080
ENV GIN_MODE=release

EXPOSE 8080

# 启动命令
CMD ["./myapp"]