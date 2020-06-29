### 构建
FROM golang:alpine as builder
WORKDIR /build
COPY . .
RUN go mod download
RUN go build -ldflags "-s -w" -o douyacun main.go

### 运行
FROM alpine as runner
WORKDIR /app
COPY --from=builder /build/douyacun /app
RUN apk add --no-cache bash
VOLUME /data
EXPOSE 9003
ENTRYPOINT ["./douyacun", "start"]