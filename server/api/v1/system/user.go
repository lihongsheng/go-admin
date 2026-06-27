package system

import (
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	"github.com/lihongsheng/go-admin/server/enum"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	"github.com/lihongsheng/go-admin/server/utils/jwt"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// UserCreate POST /system/user
func UserCreate(c *gin.Context) {
	var req dtoSys.UserCreateReq
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
	usr, err := serviceSys.DefaultUser.Create(req, mchID, sysType)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, usr)
}

// UserUpdate PUT /system/user
func UserUpdate(c *gin.Context) {
	var req dtoSys.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid body")
		return
	}
	u, err := jwt.GetUser(c.Request.Context())
	if err != nil {
		response.FailHTTP(c, 401, response.CodeUnauthorized, err.Error())
		return
	}
	mchID := u.MchID
	if u.SystemType == enum.SystemTypePlatform && req.MchID > 0 {
		mchID = req.MchID
	}
	if err := serviceSys.DefaultUser.Update(req, mchID); err != nil {
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
	u, err := jwt.GetUser(c.Request.Context())
	if err != nil {
		response.FailHTTP(c, 401, response.CodeUnauthorized, err.Error())
		return
	}
	mchID := u.MchID
	sysType := int(u.SystemType)
	if u.SystemType == enum.SystemTypePlatform {
		if req.MchID > 0 {
			mchID = req.MchID
		}
		if req.SystemType > 0 {
			sysType = req.SystemType
		} else {
			sysType = 0
		}
	}
	resp, err := serviceSys.DefaultUser.List(req, mchID, sysType)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}
