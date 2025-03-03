FROM golang:1.23.3-alpine AS builder

ARG APP_RELATIVE_PATH=./cmd/server/

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

# 设置工作目录
WORKDIR /app

COPY . .

# 设置镜像源
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod tidy && go build -ldflags="-s -w" -o server ${APP_RELATIVE_PATH}

FROM alpine:latest

# 设置工作目录
WORKDIR /app

RUN set -eux && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add tzdata && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

COPY --from=builder /app/server .
# 设置可执行权限
RUN chmod +x ./server
# 设置环境变量（如果需要）
ENV GIN_MODE=release

EXPOSE 8080
ENTRYPOINT [ "./server" ]

