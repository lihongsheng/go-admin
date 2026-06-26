package example

import (
	dtoExample "github.com/lihongsheng/go-admin/server/plugin/example/dto"
	serviceExample "github.com/lihongsheng/go-admin/server/plugin/example/service"
	"strconv"

	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

func create(c *gin.Context) {
	var req dtoExample.NoteCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, err.Error())
		return
	}
	n, err := serviceExample.DefaultNote.Create(req)
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, n)
}

func del(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		response.Fail(c, "invalid id")
		return
	}
	if err := serviceExample.DefaultNote.Delete(uint(id)); err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OKMsg(c, "ok")
}

func list(c *gin.Context) {
	ns, err := serviceExample.DefaultNote.List()
	if err != nil {
		response.Fail(c, err.Error())
		return
	}
	response.OK(c, gin.H{"list": ns, "total": len(ns)})
}
