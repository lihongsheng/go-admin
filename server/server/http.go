package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/lihongsheng/go-admin/server/log"
	"net/http"
)

// HttpServer 实现 Server 接口
type HttpServer struct {
	server *http.Server
}

// NewHttpServer 创建 HTTP 服务实例
func NewHttpServer(addr string, handler http.Handler) *HttpServer {
	return &HttpServer{
		server: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

// Start 启动 HTTP 服务（实现 Server 接口）
func (h *HttpServer) Start(ctx context.Context) error {
	log.Info("HTTP 服务启动中，监听地址：%s", h.server.Addr)
	// 使用 goroutine 启动服务，避免阻塞
	if err := h.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("HTTP 服务启动失败：%v", "err", err)
	}
	return nil
}

// Stop 优雅停止 HTTP 服务（实现 Server 接口）
func (h *HttpServer) Stop(ctx context.Context) error {
	log.Info("开始优雅停止 HTTP 服务...")
	// 调用 http.Server 的 Shutdown 方法实现优雅停止
	if err := h.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务停止失败：%w", err)
	}
	return nil
}
