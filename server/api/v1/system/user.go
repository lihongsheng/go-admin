package system

import (
	dtoSys "go-admin/server/dto/system"
	serviceSys "go-admin/server/service/system"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// UserCreate POST /system/user
func UserCreate(c *gin.Context) {
	var req dtoSys.UserCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	u, err := serviceSys.DefaultUser.Create(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, u)
}

// UserUpdate PUT /system/user
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

// UserDelete DELETE /system/user/:id
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

// UserList GET /system/user/list
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
