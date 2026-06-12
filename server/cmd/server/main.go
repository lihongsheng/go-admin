package main

import (
	"context"
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-admin/server/config"
	"go-admin/server/global"
	"go-admin/server/initialize"
	applog "go-admin/server/log"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}
	global.Cfg = cfg

	// 初始化日志系统
	logger, err := initialize.Logger(cfg.Log)
	if err != nil {
		stdlog.Fatalf("init logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck

	// 设置全局 logger
	applog.SetGlobal(logger)

	// 初始化 OpenTelemetry
	if err := initialize.InitOpenTelemetry(cfg.Observability,
		initialize.WithOTelLogger(logger),
		initialize.WithOTelAppConfig(cfg.App),
	); err != nil {
		logger.Error("init opentelemetry: " + err.Error())
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := initialize.ShutdownOpenTelemetry(ctx); err != nil {
			logger.Error("shutdown opentelemetry: " + err.Error())
		}
	}()

	// 初始化 Prometheus Metrics
	var metricsHandler http.Handler
	if err := initialize.InitMetrics(cfg.Observability.Metrics,
		initialize.WithMetricsLogger(logger),
		initialize.WithMetricsAppConfig(cfg.App),
	); err != nil {
		logger.Error("init metrics: " + err.Error())
	} else {
		if h := initialize.MetricsHandler(); h != nil {
			metricsHandler = h
		}
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := initialize.ShutdownMetrics(ctx); err != nil {
			logger.Error("shutdown metrics: " + err.Error())
		}
	}()

	// install service 早于 DB 装配——安装向导本身无需 DB 就绪
	initialize.InitInstallService()

	// 尝试连接 DB；连不上不致命，进入安装向导模式
	if cfg.DB.Configured() {
		if err := initialize.GormConnect(); err != nil {
			logger.Warn("db connect failed, entering install mode: " + err.Error())
		} else {
			initialize.DetectInstalled()
			initialize.SetupCasbin()
		}
	} else {
		logger.Warn("db not configured, entering install mode")
	}

	// 加载插件（注册 Model 到 installer 注册中心 / 注册路由）
	initialize.LoadPlugins()

	// 已安装且 DB 就绪：启动期增量同步（新插件 Model/菜单/API 自动迁入）
	initialize.SyncOnBoot()

	// 装配依赖 global.DB 的 service 单例（DB 未就绪时这步空跑，等安装回调）
	initialize.InitDBServices()

	// 启动服务器（支持优雅关闭）
	r := initialize.Router(
		initialize.WithRouterLogger(logger),
		initialize.WithRouterConfig(*cfg),
		initialize.WithMetricsHandler(metricsHandler),
	)
	addr := fmt.Sprintf(":%d", cfg.App.Port)
	logger.Info("server start",
		"addr", addr,
		"installed", global.Installed.Load(),
	)

	// 在 goroutine 中启动服务器
	srvErr := make(chan error, 1)
	go func() {
		if err := r.Run(addr); err != nil {
			srvErr <- err
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-srvErr:
		logger.Fatal("server run: " + err.Error())
	case <-quit:
		logger.Info("shutting down server...")
	}
}
