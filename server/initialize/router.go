package initialize

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"go-admin/server/config"
	"go-admin/server/log"
	"go-admin/server/middleware"
	"go-admin/server/router"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// RouterOption 路由初始化选项
type RouterOption func(*routerConfig)

type routerConfig struct {
	logger         log.Logger
	cfg            config.Config
	install        *installGuardConfig
	metricsHandler http.Handler
}

type installGuardConfig struct {
	// 这里可以放 InstallGuard 需要的参数
}

// WithRouterLogger 为路由设置 logger
func WithRouterLogger(logger log.Logger) RouterOption {
	return func(c *routerConfig) {
		c.logger = logger
	}
}

// WithRouterConfig 为路由设置配置
func WithRouterConfig(cfg config.Config) RouterOption {
	return func(c *routerConfig) {
		c.cfg = cfg
	}
}

// WithMetricsHandler 为路由设置 Prometheus metrics handler
func WithMetricsHandler(handler http.Handler) RouterOption {
	return func(c *routerConfig) {
		c.metricsHandler = handler
	}
}

// Router 构造 gin 路由
func Router(opts ...RouterOption) *gin.Engine {
	cfg := &routerConfig{
		logger: log.Global(),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	gin.SetMode(cfg.cfg.App.Mode)
	r := gin.New()
	// 中间件顺序很重要！
	r.Use(gin.Recovery())
	// 1. 先添加 otelgin 中间件，它会：
	//    - 从 header 提取 W3C Trace Context（如果有）
	//    - 或者创建新的 span（作为第一个请求）
	var serviceName string
	if cfg.cfg.Observability.Trace.Enable {
		serviceName = cfg.cfg.Observability.Trace.ServiceName
		if serviceName == "" {
			serviceName = cfg.cfg.App.Name
			if serviceName == "" {
				serviceName = "go-admin-server"
			}
		}
		r.Use(otelgin.Middleware(serviceName))
	}

	// 2. 然后添加我们的 Trace 中间件，它会：
	//    - 从 context 读取 otelgin 设置的 span
	//    - 通过 valuers 自动提取 trace_id 等信息
	r.Use(middleware.Trace(
		middleware.WithLogger(cfg.logger),
	))

	// 3. HTTP 指标中间件（记录请求计数、延迟、错误率）
	if cfg.cfg.Observability.Metrics.Enable {
		r.Use(middleware.Metrics())
	}

	// 4. 其他中间件
	r.Use(middleware.Cors())
	r.Use(middleware.RequestLog(
		middleware.WithRequestLogger(cfg.logger),
	))

	// 健康检查 / 安装状态前置（不走 InstallGuard）
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	// Prometheus metrics 端点（放在 InstallGuard 之前，确保始终可达）
	if cfg.cfg.Observability.Metrics.Enable && cfg.metricsHandler != nil {
		metricsPath := cfg.cfg.Observability.Metrics.Path
		if metricsPath == "" {
			metricsPath = "/metrics"
		}
		r.GET(metricsPath, func(c *gin.Context) {
			cfg.metricsHandler.ServeHTTP(c.Writer, c.Request)
		})
		cfg.logger.Info("metrics endpoint mounted", "path", metricsPath)
	}

	// 上传文件静态服务（与配置中的 upload.local.path 对应）
	// 仅当使用本地存储时生效；云存储文件通过 CDN/OSS 直接访问
	uploadPath := cfg.cfg.Upload.Local.Path
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	// 图片文件加载接口 —— 支持长缓存（文件名含 MD5 哈希，内容变更即文件名变更）
	r.GET("/uploads/*filepath", func(c *gin.Context) {
		fp := c.Param("filepath")
		if strings.Contains(fp, "..") {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		ext := filepath.Ext(fp)
		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp", ".ico":
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.File(filepath.Join(uploadPath, fp))
	})

	// 安装向导直接挂在根路径（无 /api 前缀，方便前端独立模块）
	router.InstallRouter(&r.RouterGroup)

	// 业务路由统一在 /api/v1 下，由 InstallGuard 拦截
	api := r.Group("/api/v1", middleware.InstallGuard())
	router.BaseRouter(api)   // 登录 / 当前用户 / 当前菜单
	router.SystemRouter(api) // user / role / menu / api
	router.PluginRouter(r)   // 已装插件列表 + 插件自身路由
	return r
}

// RequestLog 修复一下 RequestLog，需要 time 包
// 让我们更新 middleware 包中的实现
func RequestLog(logger log.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		reqLogger := log.FromContext(c.Request.Context())
		reqLogger.Info("http",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"cost", time.Since(start),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)
	}
}
