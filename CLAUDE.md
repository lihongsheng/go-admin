
# CLAUDE.md

该文件可为 Claude Code（claude.ai/code）在本代码仓库中处理代码时提供操作指引。

## 项目概述

这是一款 基于Go语言和vue 以及 Element 编写的一个web后台管理系统，包含用户管理，角色管理，菜单管理，api管理。菜单支持vue按钮权限管控。支持插件化新增页面，插件化导入新的后端api和菜单栏。支持otel和普罗米修斯监控。
支持多商户，商户管理员具有 商户用户，商户角色管理，不具有商户菜单和API维护权限。
## 编码规范

本项目遵循 Go 高级开发实践，请注意以下原则：

### 依赖注入
- 所有 service 层通过构造函数注入依赖（如 `NewXxxService(repo, casbin, ...)`），不允许在内部直接 `new` 或使用全局变量拼凑
- repo / casbin / 其他外部依赖统一在 `initialize/service.go` 中装配好再注入
- 包级单例（`DefaultXxx`）只用于 wiring，业务代码不直接引用包级单例之外的全局变量

### 接口编程
- service 层必须定义接口（如 `type XxxService interface { ... }`），返回接口类型，让调用方面向接口编程
- repo 层同样定义接口，方便单测替换
- casbin 操作通过 `Port` 接口暴露，不直接导入 `casbin` 包

### 分层架构
- handler（`api/v1/`）→ service（`service/`）→ repo（`repo/`），不允许跨层调用
- handler 只做参数校验和响应组装，不写业务逻辑
- middleware 不直接操作 DB，通过 service 方法完成业务检查
- 跨层传递用 DTO（`dto/`），不直接暴露 model

### 不要硬编码
- 枚举值用 `enum/` 目录下的类型常量，不写 magic number
- 不要用 `"u:" + id` 这种前缀拼接 — casbin 策略直接使用数字 ID
- 不要用角色 code 做策略主体 — 一律用角色 ID
- 配置项走 `global.Cfg`，不硬编码在业务代码中

## Commands

## 架构

### server

server目录是服务器代码的根目录，基于go gin 框架开发，基于插件化开发，支持插件化新增页面，插件化导入新的后端api和菜单栏。

### web 目录
web 目录是vue 前端代码的根目录，基于element 框架开发, 基于插件化开发，支持插件化新增页面。
