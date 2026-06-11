package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 服务端总配置
type Config struct {
	App           App           `mapstructure:"app"           json:"app"           yaml:"app"`
	JWT           JWT           `mapstructure:"jwt"           json:"jwt"           yaml:"jwt"`
	Log           Log           `mapstructure:"log"           json:"log"           yaml:"log"`
	DB            DB            `mapstructure:"db"            json:"db"            yaml:"db"`
	Redis         Redis         `mapstructure:"redis"         json:"redis"         yaml:"redis"`
	Install       Install       `mapstructure:"install"       json:"install"       yaml:"install"`
	Observability Observability `mapstructure:"observability" json:"observability" yaml:"observability"`
}

type App struct {
	Name string `mapstructure:"name" json:"name" yaml:"name"`
	Mode string `mapstructure:"mode" json:"mode" yaml:"mode"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}

type JWT struct {
	Secret string `mapstructure:"secret" json:"secret" yaml:"secret"`
	Expire int64  `mapstructure:"expire" json:"expire" yaml:"expire"`
	Issuer string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
}

type Log struct {
	Level      string `mapstructure:"level"       json:"level"       yaml:"level"`
	Dir        string `mapstructure:"dir"         json:"dir"         yaml:"dir"`
	MaxSize    int    `mapstructure:"max_size"    json:"max_size"    yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"     json:"max_age"     yaml:"max_age"`
}

// DB 数据库配置；driver 为空时视为未配置
type DB struct {
	Driver   string `mapstructure:"driver"   json:"driver"   yaml:"driver"`
	Host     string `mapstructure:"host"     json:"host"     yaml:"host"`
	Port     int    `mapstructure:"port"     json:"port"     yaml:"port"`
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Database string `mapstructure:"database" json:"database" yaml:"database"`
	Path     string `mapstructure:"path"     json:"path"     yaml:"path"` // sqlite
	Charset  string `mapstructure:"charset"  json:"charset"  yaml:"charset"`
	MaxIdle  int    `mapstructure:"max_idle" json:"max_idle" yaml:"max_idle"`
	MaxOpen  int    `mapstructure:"max_open" json:"max_open" yaml:"max_open"`
	LogMode  string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
}

// Configured 判断 DB 是否已配置
func (d DB) Configured() bool { return d.Driver != "" }

// DSN 根据 driver 拼装连接串（连接目标数据库）
func (d DB) DSN() string {
	switch d.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Port, d.Database, d.Charset)
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			d.Host, d.Username, d.Password, d.Database, d.Port)
	case "sqlite":
		if d.Path == "" {
			return "go-admin.db"
		}
		return d.Path
	}
	return ""
}

// ServerDSN 返回不带具体 database 的连接串，用于 CREATE DATABASE
// MySQL 连到默认实例；Postgres 连到 "postgres" 系统库。SQLite 无意义返回空串。
func (d DB) ServerDSN() string {
	switch d.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
			d.Username, d.Password, d.Host, d.Port, d.Charset)
	case "postgres":
		return fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%d sslmode=disable TimeZone=Asia/Shanghai",
			d.Host, d.Username, d.Password, d.Port)
	}
	return ""
}

type Redis struct {
	Enable   bool   `mapstructure:"enable"   json:"enable"   yaml:"enable"`
	Addr     string `mapstructure:"addr"     json:"addr"     yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db"       json:"db"       yaml:"db"`
}

type Install struct {
	Enable bool `mapstructure:"enable" json:"enable" yaml:"enable"`
}

// Observability 可观测性配置
type Observability struct {
	// Trace 分布式追踪配置
	Trace Trace `mapstructure:"trace" json:"trace" yaml:"trace"`
	// Metrics 指标配置
	Metrics Metrics `mapstructure:"metrics" json:"metrics" yaml:"metrics"`
}

// Trace 分布式追踪配置
type Trace struct {
	// Enable 是否启用追踪
	Enable bool `mapstructure:"enable" json:"enable" yaml:"enable"`
	// Exporter 导出器类型: stdout, otlp
	Exporter string `mapstructure:"exporter" json:"exporter" yaml:"exporter"`
	// Endpoint OTLP 接收端地址，例如 "localhost:4317"(gRPC) 或 "http://localhost:4318"(HTTP)
	// 仅当 Exporter 为 "otlp" 时生效
	Endpoint string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	// ServiceName 服务名称
	ServiceName string `mapstructure:"service_name" json:"service_name" yaml:"service_name"`
	// SampleRate 采样率 0.0-1.0
	SampleRate float64 `mapstructure:"sample_rate" json:"sample_rate" yaml:"sample_rate"`
}

// Metrics 指标配置
type Metrics struct {
	// Enable 是否启用指标
	Enable bool `mapstructure:"enable" json:"enable" yaml:"enable"`
	// Exporter 导出器类型: prometheus（预留扩展）
	Exporter string `mapstructure:"exporter" json:"exporter" yaml:"exporter"`
	// Endpoint OTLP metrics 推送地址，例如 "localhost:4317"(gRPC) 或 "http://localhost:4318"(HTTP)
	// 设置后将同时推送指标到 OTLP collector（Prometheus /metrics 端点不受影响）
	Endpoint string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	// Path 指标暴露路径
	Path string `mapstructure:"path" json:"path" yaml:"path"`
}

var (
	V       *viper.Viper
	cfgPath string
)

// Load 读取 config.yaml
func Load(path string) (*Config, error) {
	cfgPath = path
	V = viper.New()
	V.SetConfigFile(path)
	V.SetConfigType("yaml")
	if err := V.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	c := &Config{}
	if err := V.Unmarshal(c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}
	return c, nil
}

// Path 返回当前配置文件路径
func Path() string { return cfgPath }

// Save 写回 config.yaml（安装向导提交 DB 配置时调用）
func Save(c *Config) error {
	V.Set("app", c.App)
	V.Set("jwt", c.JWT)
	V.Set("log", c.Log)
	V.Set("db", c.DB)
	V.Set("redis", c.Redis)
	V.Set("install", c.Install)
	V.Set("observability", c.Observability)
	return V.WriteConfig()
}
