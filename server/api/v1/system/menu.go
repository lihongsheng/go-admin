package system

import (
	modelSys "github.com/lihongsheng/go-admin/server/model/system"
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// swagger type hints
var (
	_ = (*modelSys.SysMenu)(nil)
	_ = (*dtoSys.MenuTreeResp)(nil)
)

// MenuCreate 新增菜单
// @Summary      新增菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.MenuCreateReq  true  "菜单信息"
// @Success      200   {object}  response.Body{data=modelSys.SysMenu}
// @Security     BearerAuth
// @Router       /api/v1/system/menu [post]
func MenuCreate(c *gin.Context) {
	var req dtoSys.MenuCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	m, err := serviceSys.DefaultMenu.Create(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, m)
}

// MenuUpdate 更新菜单
// @Summary      更新菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Param        body  body      dtoSys.MenuUpdateReq  true  "菜单信息"
// @Success      200   {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/menu [put]
func MenuUpdate(c *gin.Context) {
	var req dtoSys.MenuUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid body")
		return
	}
	if err := serviceSys.DefaultMenu.Update(req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// MenuDelete 递归删除菜单节点及全部子孙
// @Summary      删除菜单（递归删除子节点）
// @Tags         菜单管理
// @Produce      json
// @Param        id   path      int  true  "菜单ID"
// @Success      200  {object}  response.Body
// @Security     BearerAuth
// @Router       /api/v1/system/menu/{id} [delete]
func MenuDelete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if err := serviceSys.DefaultMenu.Delete(id); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

// MenuTree 全量菜单树
// @Summary      获取全量菜单树
// @Tags         菜单管理
// @Produce      json
// @Success      200  {object}  response.Body{data=dtoSys.MenuTreeResp}
// @Security     BearerAuth
// @Router       /api/v1/system/menu/tree [get]
func MenuTree(c *gin.Context) {
	resp, err := serviceSys.DefaultMenu.Tree()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}
