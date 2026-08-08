package system

import (
	modelSys "github.com/lihongsheng/go-admin/server/model/system"
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// swagger type hints — 让 swag 能解析 @Success 中的 modelSys.XXX / dtoSys.XXX 类型引用
var (
	_ = (*modelSys.SysUser)(nil)
	_ = (*dtoSys.UserListResp)(nil)
)

// UserCreate 新增用户
// @Summary      新增用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.UserCreateReq  true  "用户信息"
// @Success      200   {object}  response.Body{data=modelSys.SysUser}
// @Security     BearerAuth
// @Router       /api/v1/system/user [post]
func UserCreate(c *gin.Context) {
	var req dtoSys.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	usr, err := serviceSys.DefaultUser.Create(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, usr)
}

// UserUpdate 更新用户
// @Summary      更新用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.UserUpdateReq  true  "用户信息"
// @Success      200   {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/user [put]
func UserUpdate(c *gin.Context) {
	var req dtoSys.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid body")
		return
	}
	if err := serviceSys.DefaultUser.Update(req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// UserDelete 删除用户
// @Summary      删除用户
// @Tags         用户管理
// @Produce      json
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/user/{id} [delete]
func UserDelete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := serviceSys.DefaultUser.Delete(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// UserList 用户列表
// @Summary      用户列表
// @Tags         用户管理
// @Produce      json
// @Param        keyword  query     string  false  "搜索关键词"
// @Param        page     query     int     false  "页码"
// @Param        size     query     int     false  "每页条数"
// @Success      200      {object}  response.Body{data=dtoSys.UserListResp}
// @Security     BearerAuth
// @Router       /api/v1/system/user/list [get]
func UserList(c *gin.Context) {
	var req dtoSys.UserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	resp, err := serviceSys.DefaultUser.List(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}
