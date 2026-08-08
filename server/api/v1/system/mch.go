package system

import (
	"strconv"

	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	"github.com/lihongsheng/go-admin/server/enum"
	modelSys "github.com/lihongsheng/go-admin/server/model/system"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// swagger type hints
var _ = (*modelSys.Merchant)(nil)

// MchCreate 新增商户
// @Summary      新增商户
// @Tags         商户管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.MchCreateRequest  true  "商户信息"
// @Success      200   {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/mch [post]
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

// MchUpdate 更新商户（与创建共用结构体，ID 决定新增或更新）
// @Summary      更新商户
// @Tags         商户管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.MchCreateRequest  true  "商户信息（需包含 ID）"
// @Success      200   {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/mch [put]
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

// MchDetail 商户详情
// @Summary      根据 ID 获取商户详情
// @Tags         商户管理
// @Produce      json
// @Param        id   path      int64  true  "商户ID"
// @Success      200  {object}  response.Body{data=modelSys.Merchant}
// @Security     BearerAuth
// @Router       /api/v1/system/mch/{id} [get]
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

// MchDetailByNo 根据编号获取商户详情
// @Summary      根据编号获取商户详情
// @Tags         商户管理
// @Produce      json
// @Param        mchNo  path      string  true  "商户编号"
// @Success      200    {object}  response.Body{data=modelSys.Merchant}
// @Security     BearerAuth
// @Router       /api/v1/system/mch/no/{mchNo} [get]
func MchDetailByNo(c *gin.Context) {
	mchNo := c.Param("mchNo")
	mch, err := serviceSys.DefaultMch.GetByMchNo(c.Request.Context(), mchNo)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, mch)
}

// MchList 商户列表
// @Summary      商户列表
// @Tags         商户管理
// @Produce      json
// @Param        mch_name  query     string  false  "公司名称"
// @Param        mch_no    query     string  false  "商户编号"
// @Param        status    query     int     false  "状态（1正常 2停用）"
// @Param        limit     query     int     false  "每页条数（最大50）"
// @Param        page      query     int     false  "页码"
// @Success      200       {object}  response.Body{data=object{list=array,total=integer}}
// @Security     BearerAuth
// @Router       /api/v1/system/mch/list [get]
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

// MchChangeStatus 修改商户状态（暂未注册到路由，保留以备后用）
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
