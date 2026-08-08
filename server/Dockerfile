# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 国内构建可取消下行注释，配置 GOPROXY 加速
# ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 生成 Swagger 文档（依赖 handler 注解），然后编译
RUN go run github.com/swaggo/swag/cmd/swag@v1.16.6 init -g cmd/server/main.go -o docs --parseDependency --parseInternal
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server/main.go

# ---- Run Stage ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config/config.docker.yaml ./config/config.yaml

RUN mkdir -p /app/logs /app/uploads

EXPOSE 8989

ENTRYPOINT ["./server", "-c", "config/config.yaml"]
