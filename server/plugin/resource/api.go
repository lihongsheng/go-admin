// Package resource 资源管理插件：云资源（ECS）实例管理
package resource

import (
	"strconv"

	dtoResource "go-admin/server/plugin/resource/dto"
	serviceResource "go-admin/server/plugin/resource/service"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// create 新建云资源
// @Summary      新增资源
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        body body dtoResource.ResourceCreateReq true "资源信息"
// @Success      200  {object}  response.Body
// @Security     BearerAuth
// @Router       /api/plugin/resource/v1 [post]
func create(c *gin.Context) {
	var req dtoResource.ResourceCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	res, err := serviceResource.DefaultResource.Create(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, res)
}

// update 更新云资源
// @Summary      编辑资源
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        id   path int true "资源 ID"
// @Param        body body dtoResource.ResourceUpdateReq true "资源信息"
// @Success      200  {object}  response.Body
// @Security     BearerAuth
// @Router       /api/plugin/resource/v1/{id} [put]
func update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, "invalid id")
		return
	}
	var req dtoResource.ResourceUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	req.ID = uint(id)
	res, err := serviceResource.DefaultResource.Update(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, res)
}

// detail 云资源详情
// @Summary      资源详情
// @Tags         资源管理
// @Produce      json
// @Param        id path int true "资源 ID"
// @Success      200   {object}  response.Body
// @Security     BearerAuth
// @Router       /api/plugin/resource/v1/{id} [get]
func detail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, "invalid id")
		return
	}
	res, err := serviceResource.DefaultResource.Get(uint(id))
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, res)
}

// batchDel 批量删除云资源
// @Summary      批量删除资源
// @Tags         资源管理
// @Accept       json
// @Produce      json
// @Param        body body dtoResource.ResourceBatchDeleteReq true "资源 ID 列表"
// @Success      200 {object} response.Body
// @Security     BearerAuth
// @Router       /api/plugin/resource/v1/batch [delete]
func batchDel(c *gin.Context) {
	var req dtoResource.ResourceBatchDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	if err := serviceResource.DefaultResource.Delete(req.IDs); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// list 云资源分页列表
// @Summary      资源列表
// @Tags         资源管理
// @Produce      json
// @Param        page    query int    false "页码"     default(1)
// @Param        limit   query int    false "每页数量" default(10)
// @Param        keyword query string false "资源名称关键字"
// @Success      200   {object}  response.Body{data=object{list=array,total=int64}}
// @Security     BearerAuth
// @Router       /api/plugin/resource/v1/list [get]
func list(c *gin.Context) {
	var req dtoResource.ResourceListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	resources, total, err := serviceResource.DefaultResource.List(req.Page, req.Limit, req.Keyword)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, gin.H{"list": resources, "total": total})
}
