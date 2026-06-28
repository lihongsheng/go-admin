// Package base 基础接口：登录 / 当前用户 / 当前菜单
package base

import (
  dtoBase "github.com/lihongsheng/go-admin/server/dto/base"
  serviceBase "github.com/lihongsheng/go-admin/server/service/base"
  serviceSys "github.com/lihongsheng/go-admin/server/service/system"
  "github.com/lihongsheng/go-admin/server/utils/jwt"
  "github.com/lihongsheng/go-admin/server/utils/response"

  "github.com/gin-gonic/gin"
)

// Captcha GET /api/v1/base/captcha 生成图形验证码
func Captcha(c *gin.Context) {
  resp, err := serviceBase.Default.Captcha()
  if err != nil {
    response.Fail(c, err.Error())
    return
  }
  response.OK(c, resp)
}

// Login POST /api/v1/base/login
func Login(c *gin.Context) {
  var req dtoBase.LoginReq
  if err := c.ShouldBindJSON(&req); err != nil {
    response.Fail(c, err.Error())
    return
  }
  resp, err := serviceBase.Default.Login(req)
  if err != nil {
    response.Fail(c, err.Error())
    return
  }
  response.OK(c, resp)
}

// Logout POST /api/v1/base/logout —— 简单版（前端清 token 即可）
func Logout(c *gin.Context) { response.OKMsg(c, "ok") }

// Info GET /api/v1/base/info —— 当前用户信息
func Info(c *gin.Context) {
  user, err := jwt.GetUser(c.Request.Context())
  if err != nil {
    response.FailHTTP(c, 401, response.CodeUnauthorized, err.Error())
    return
  }
  u, err := serviceBase.Default.Info(user.ID)
  if err != nil {
    response.Fail(c, err.Error())
    return
  }
  response.OK(c, u)
}

// Menu GET /api/v1/base/menu —— 当前用户菜单（树）
func Menu(c *gin.Context) {
  user, err := jwt.GetUser(c.Request.Context())
  if err != nil {
    response.FailHTTP(c, 401, response.CodeUnauthorized, err.Error())
    return
  }
  menus, err := serviceSys.DefaultMenu.UserTree(user.ID)
  if err != nil {
    response.Fail(c, err.Error())
    return
  }
  response.OK(c, gin.H{"menus": menus})
}
