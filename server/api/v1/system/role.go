package system

import (
	dtoSys "go-admin/server/dto/system"
	serviceSys "go-admin/server/service/system"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// RoleCreate POST /system/role
func RoleCreate(c *gin.Context) {
	var req dtoSys.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	r, err := serviceSys.DefaultRole.Create(req)
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
	resp, err := serviceSys.DefaultRole.List()
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
