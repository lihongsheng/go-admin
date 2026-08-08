// @title           go-admin API
// @version         1.0.0
// @description     基于 Gin + Vue3 的后台管理系统 API。<br/>统一响应格式：<code>{"code":0,"msg":"ok","data":...}</code>，code=0 成功，code=7 业务失败。<br/>所有 <code>/api/v1/system/*</code> 接口需要 <code>Authorization: Bearer &lt;token&gt;</code>。
// @host            localhost:8989
// @BasePath        /
// @securityDefinitions.apikey  BearerAuth
// @in              header
// @name            Authorization
package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lihongsheng/go-admin/server/cron"
	"github.com/lihongsheng/go-admin/server/cron/initalize"
	"github.com/lihongsheng/go-admin/server/server"
	stdlog "log"
	"net/http"
	"os"
	"time"

	"github.com/lihongsheng/go-admin/server/config"
	_ "github.com/lihongsheng/go-admin/server/docs" // swag 生成的 OpenAPI 文档
	"github.com/lihongsheng/go-admin/server/global"
	"github.com/lihongsheng/go-admin/server/initialize"
	applog "github.com/lihongsheng/go-admin/server/log"
)

func main() {
	cfgPath := flag.String("c", "config/config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		stdlog.Fatalf("load config: %v", err)
	}
	global.Cfg = cfg

	// 根据配置文件端口动态设置 swagger host
	// docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.App.Port)

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

	// 初始化 Redis（如启用）
	initialize.InitRedis()
	// 初始化验证码存储（根据配置选择 memory / redis）
	initialize.InitCaptcha()
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

	// 启动服务
	httpSrv := server.NewHttpServer(addr, r)
	cronSrv := cron.NewCronServer(initalize.GetCronJobs()...)
	app := server.NewApp(30*time.Second, httpSrv, cronSrv)
	if err := app.Run(); err != nil {
		logger.Error("server run failed: " + err.Error())
		os.Exit(1)
	}
	logger.Info("server stopped")
}
