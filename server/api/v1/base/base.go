// Package base 基础接口：登录 / 当前用户 / 当前菜单
package base

import (
	dtoBase "go-admin/server/dto/base"
	modelSys "go-admin/server/model/system"
	serviceBase "go-admin/server/service/base"
	serviceSys "go-admin/server/service/system"
	"go-admin/server/utils/jwt"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// swagger type hints
var (
	_ = (*modelSys.SysUser)(nil)
	_ = (*dtoBase.CaptchaResp)(nil)
	_ = (*dtoBase.LoginResp)(nil)
)

// Captcha 获取图形验证码
// @Summary      获取图形验证码
// @Tags         基础接口
// @Produce      json
// @Success      200  {object}  response.Body{data=dtoBase.CaptchaResp}
// @Router       /api/v1/base/captcha [get]
func Captcha(c *gin.Context) {
	resp, err := serviceBase.Default.Captcha()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// Login 用户登录
// @Summary      用户登录
// @Tags         基础接口
// @Accept       json
// @Produce      json
// @Param        body  body      dtoBase.LoginReq  true  "登录参数"
// @Success      200   {object}  response.Body{data=dtoBase.LoginResp}
// @Router       /api/v1/base/login [post]
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

// Logout 退出登录
// @Summary      退出登录
// @Tags         基础接口
// @Produce      json
// @Success      200  {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/base/logout [post]
func Logout(c *gin.Context) { response.OKMsg(c, "ok") }

// Info 当前用户信息
// @Summary      获取当前登录用户信息
// @Tags         基础接口
// @Produce      json
// @Success      200  {object}  response.Body{data=modelSys.SysUser}
// @Failure      401  {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/base/info [get]
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

// Menu 当前用户菜单树
// @Summary      获取当前用户菜单树
// @Tags         基础接口
// @Produce      json
// @Success      200  {object}  response.Body{data=object{menus=array}}
// @Failure      401  {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/base/menu [get]
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
