package system

import (
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	serviceSys "github.com/lihongsheng/go-admin/server/service/system"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// MenuCreate POST /system/menu
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

// MenuUpdate PUT /system/menu
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

// MenuDelete DELETE /system/menu/:id —— 递归删除目标节点及全部子孙
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

// MenuTree GET /system/menu/tree —— 全量菜单树
func MenuTree(c *gin.Context) {
	resp, err := serviceSys.DefaultMenu.Tree()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, resp)
}
