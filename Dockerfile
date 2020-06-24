### 构建
FROM golang:alpine as builder
WORKDIR /build
ENV GOPROXY=https://goproxy.cn
COPY . .
RUN go mod download && go build -ldflags "-s -w" -a -o douyacun main.go

### 运行
FROM Alpine as runner
COPY --from=builder /build/douyacun /app
CMD douyacun start --env prod
