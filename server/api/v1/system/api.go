package system

import (
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	"github.com/lihongsheng/go-admin/server/utils/jwt"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// ApiCreate POST /system/api
func ApiCreate(c *gin.Context) {
	var req dtoSys.ApiCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	// 如果前端没有传 system_type，从 JWT 中取（商户管理员）
	if req.SystemType == 0 {
		if u, err := jwt.GetUser(c.Request.Context()); err == nil {
			req.SystemType = u.SystemType
		}
	}
	a, err := serviceSys.DefaultApi.Create(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, a)
}

// ApiUpdate PUT /system/api
func ApiUpdate(c *gin.Context) {
	var req dtoSys.ApiUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid body")
		return
	}
	if err := serviceSys.DefaultApi.Update(req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// ApiDelete DELETE /system/api/:id
func ApiDelete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := serviceSys.DefaultApi.Delete(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// ApiList GET /system/api/list
// 支持 ?system_type=N 按系统类型过滤
func ApiList(c *gin.Context) {
	var req dtoSys.ApiListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	resp, err := serviceSys.DefaultApi.List(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}
