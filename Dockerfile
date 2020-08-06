### 构建
FROM golang:alpine as builder
WORKDIR /build
COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && go mod download
RUN go build -ldflags "-s -w" -o douyacun main.go

### 运行
FROM registry.cn-hangzhou.aliyuncs.com/douyacun/alpine:latest as runner
WORKDIR /app
COPY --from=builder /build/douyacun /bin
VOLUME /data
EXPOSE 9003
ENTRYPOINT ["douyacun", "start"]