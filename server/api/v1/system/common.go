// Package system system 模块 HTTP handler（只做参数校验 + service 调用 + 响应包装）
package system

import (
	"strconv"

	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// parseID 从 URL :id 提取 uint；不合法时直接写入失败响应并返回 false
func parseID(c *gin.Context) (uint, bool) {
	v, err := strconv.Atoi(c.Param("id"))
	if err != nil || v <= 0 {
		response.Fail(c, "invalid id")
		return 0, false
	}
	return uint(v), true
}
