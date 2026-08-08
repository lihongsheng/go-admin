# go-admin 项目 Makefile
SHELL := /bin/bash
SERVER_DIR := server
SWAG_VERSION := v1.16.6
SWAG_CMD := swag

.PHONY: swagger swagger-fmt run dev build test tidy install clean

swag:
	go install github.com/swaggo/swag/cmd/swag@v1.16.6

## 生成 Swagger 文档到 server/docs/
swagger:
	cd $(SERVER_DIR) && $(SWAG_CMD) init -g cmd/server/main.go -o docs --parseDependency --parseInternal

## 格式化 swagger 注解
swagger-fmt:
	cd $(SERVER_DIR) && $(SWAG_CMD) fmt -g cmd/server/main.go

## 启动后端（开发模式）
run:
	cd $(SERVER_DIR) && go run cmd/server/main.go -c config/config.yaml

## 启动后端（air 热重载，需已安装 air）
dev:
	cd $(SERVER_DIR) && air

## 编译后端
build:
	cd $(SERVER_DIR) && go build -ldflags="-s -w" -o bin/server ./cmd/server/main.go

## 运行测试
test:
	cd $(SERVER_DIR) && go test ./...

## 整理依赖
tidy:
	cd $(SERVER_DIR) && go mod tidy

## 下载依赖
install:
	cd $(SERVER_DIR) && go mod download

## 清理构建产物
clean:
	rm -rf $(SERVER_DIR)/bin
