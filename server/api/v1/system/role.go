package system

import (
	dtoSys "go-admin/server/dto/system"
	modelSys "go-admin/server/model/system"
	serviceSys "go-admin/server/service/system"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// swagger type hints
var (
	_ = (*modelSys.SysRole)(nil)
	_ = (*dtoSys.RoleListResp)(nil)
	_ = (*dtoSys.RoleAuthDetailResp)(nil)
)

// RoleCreate 新增角色
// @Summary      新增角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.RoleCreateReq  true  "角色信息"
// @Success      200   {object}  response.Body{data=modelSys.SysRole}
// @Security     BearerAuth
// @Router       /api/v1/system/role [post]
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

// RoleUpdate 更新角色
// @Summary      更新角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.RoleUpdateReq  true  "角色信息"
// @Success      200   {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/role [put]
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

// RoleDelete 删除角色
// @Summary      删除角色
// @Tags         角色管理
// @Produce      json
// @Param        id   path      int  true  "角色ID"
// @Success      200  {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/role/{id} [delete]
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

// RoleList 角色列表
// @Summary      角色列表
// @Tags         角色管理
// @Produce      json
// @Success      200  {object}  response.Body{data=dtoSys.RoleListResp}
// @Security     BearerAuth
// @Router       /api/v1/system/role/list [get]
func RoleList(c *gin.Context) {
	var req dtoSys.RoleListReq
	_ = c.ShouldBindQuery(&req)
	resp, err := serviceSys.DefaultRole.List(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}

// RoleAuth 设置角色的菜单/API 权限
// @Summary      设置角色的菜单/API 权限
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.RoleAuthReq  true  "授权参数"
// @Success      200   {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/role/auth [post]
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

// RoleAuthDetail 查询角色已分配的菜单/API
// @Summary      查询角色已分配的菜单/API
// @Tags         角色管理
// @Produce      json
// @Param        id   path      int  true  "角色ID"
// @Success      200  {object}  response.Body{data=dtoSys.RoleAuthDetailResp}
// @Security     BearerAuth
// @Router       /api/v1/system/role/auth/{id} [get]
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

// RoleSetDefaultRouter 设置角色默认首页路由
// @Summary      设置角色默认首页路由
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Param        id    path      int                          true  "角色ID"
// @Param        body  body      dtoSys.RoleSetDefaultRouterReq  true  "默认路由"
// @Success      200   {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/role/{id}/default-router [put]
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
