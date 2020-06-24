### 构建
FROM golang:alpine as builder
WORKDIR /build
ENV GOPROXY=https://goproxy.cn
COPY . .
RUN go mod download; \
    go build -ldflags "-s -w" -o douyacun main.go

### 运行
FROM alpine as runner
COPY --from=builder /build/douyacun /app
VOLUME /data
EXPOSE 9003
CMD douyacun start --env prod
