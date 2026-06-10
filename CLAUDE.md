
# CLAUDE.md

该文件可为 Claude Code（claude.ai/code）在本代码仓库中处理代码时提供操作指引。

## 项目概述

这是一款 基于Go语言和vue 以及 Element 编写的一个web后台管理系统，包含用户管理，角色管理，菜单管理，api管理。菜单支持vue按钮权限管控。项目架构类似于go-vue-admin 和go-admin,支持插件化新增页面，插件化导入新的后端api和菜单栏。菜单栏管控参考go-admin,返回给vue符合vue渲染菜单栏规范。
## Commands

## 架构

### server

server目录是服务器代码的根目录，基于go gin 框架开发，基于插件化开发，支持插件化新增页面，插件化导入新的后端api和菜单栏。

### web 目录
web 目录是vue 前端代码的根目录，基于element 框架开发, 基于插件化开发，支持插件化新增页面。