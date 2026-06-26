package system

import (
	"strconv"

	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	"github.com/lihongsheng/go-admin/server/enum"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// MchCreate POST /system/mch
func MchCreate(c *gin.Context) {
	var req dtoSys.MchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err := serviceSys.DefaultMch.Create(c.Request.Context(), req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// MchUpdate PUT /system/mch
func MchUpdate(c *gin.Context) {
	var req dtoSys.MchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid body")
		return
	}
	if err := req.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err := serviceSys.DefaultMch.Save(c.Request.Context(), req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// MchDetail GET /system/mch/:id
func MchDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, "invalid id")
		return
	}
	mch, err := serviceSys.DefaultMch.Get(c.Request.Context(), id)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, mch)
}

// MchDetailByNo GET /system/mch/no/:mchNo
func MchDetailByNo(c *gin.Context) {
	mchNo := c.Param("mchNo")
	mch, err := serviceSys.DefaultMch.GetByMchNo(c.Request.Context(), mchNo)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, mch)
}

// MchList GET /system/mch/list
func MchList(c *gin.Context) {
	var req dtoSys.MchQueryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err := req.Validate(); err != nil {
		response.Fail(c, err.Error())
		return
	}
	list, err := serviceSys.DefaultMch.Search(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	total, err := serviceSys.DefaultMch.Count(c.Request.Context(), req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, gin.H{"list": list, "total": total})
}

// MchChangeStatus PUT /system/mch/status
func MchChangeStatus(c *gin.Context) {
	var req dtoSys.MchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if req.ID <= 0 {
		response.Fail(c, "商户ID不能为空")
		return
	}
	if req.Status < 1 {
		response.Fail(c, "商户状态不能为空")
		return
	}
	if err := serviceSys.DefaultMch.ChangeStatus(c.Request.Context(), req.ID, enum.MchStatus(req.Status)); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}
