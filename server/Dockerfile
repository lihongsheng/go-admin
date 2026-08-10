# ---- Build Stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 国内构建可取消下行注释，配置 GOPROXY 加速
ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build  -o ser ./cmd/server/main.go

# ---- Run Stage ----
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/ser .
COPY --from=builder /app/config/config.docker.yaml ./config/config.yaml

RUN mkdir -p /app/logs /app/uploads

EXPOSE 8989

ENTRYPOINT ["./ser", "-c", "config/config.yaml"]
