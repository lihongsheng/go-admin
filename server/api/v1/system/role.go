package system

import (
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	"github.com/lihongsheng/go-admin/server/enum"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	"github.com/lihongsheng/go-admin/server/utils/jwt"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// RoleCreate POST /system/role
func RoleCreate(c *gin.Context) {
	var req dtoSys.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	u, err := jwt.GetUser(c.Request.Context())
	if err != nil {
		response.FailHTTP(c, 401, response.CodeUnauthorized, err.Error())
		return
	}
	mchID := u.MchID
	sysType := u.SystemType
	if u.SystemType == enum.SystemTypePlatform {
		if req.MchID > 0 {
			mchID = req.MchID
		}
		if req.SystemType > 0 {
			sysType = req.SystemType
		}
	}
	r, err := serviceSys.DefaultRole.Create(req, mchID, sysType)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, r)
}

// RoleUpdate PUT /system/role
func RoleUpdate(c *gin.Context) {
	var req dtoSys.RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid body")
		return
	}
	if err := serviceSys.DefaultRole.Update(req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// RoleDelete DELETE /system/role/:id
func RoleDelete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := serviceSys.DefaultRole.Delete(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// RoleList GET /system/role/list
func RoleList(c *gin.Context) {
	u, err := jwt.GetUser(c.Request.Context())
	if err != nil {
		response.FailHTTP(c, 401, response.CodeUnauthorized, err.Error())
		return
	}
	var req dtoSys.RoleListReq
	_ = c.ShouldBindQuery(&req)
	if req.SystemType == 0 {
		req.SystemType = int(u.SystemType)
	}
	if req.MchID == 0 {
		req.MchID = u.MchID
	}
	resp, err := serviceSys.DefaultRole.List(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// RoleAuth POST /system/role/auth —— 设置角色的菜单/API 权限
func RoleAuth(c *gin.Context) {
	var req dtoSys.RoleAuthReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err := serviceSys.DefaultRole.Auth(req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// RoleAuthDetail GET /system/role/auth/:id —— 查询角色已分配的菜单/API
func RoleAuthDetail(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	resp, err := serviceSys.DefaultRole.AuthDetail(id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// RoleSetDefaultRouter PUT /system/role/:id/default-router —— 设置角色默认首页路由
func RoleSetDefaultRouter(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dtoSys.RoleSetDefaultRouterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err := serviceSys.DefaultRole.SetDefaultRouter(id, req.DefaultRouter); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}
