package middleware

import (
	"strings"

	"github.com/lihongsheng/go-admin/server/global"
	"github.com/lihongsheng/go-admin/server/utils/jwt"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth 校验 Authorization: Bearer xxx
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		raw := c.GetHeader("Authorization")
		if raw == "" {
			raw = c.Query("token")
		}
		raw = strings.TrimPrefix(raw, "Bearer ")
		if raw == "" {
			response.FailCode(c, response.CodeUnauthorized, "missing token")
			c.Abort()
			return
		}
		claim, err := jwt.Parse(raw, global.Cfg.JWT)
		if err != nil {
			response.FailCode(c, response.CodeUnauthorized, "invalid token: "+err.Error())
			c.Abort()
			return
		}
		// 将用户信息写入请求上下文，供后续中间件 / handler 通过 jwt.GetUser() 获取
		ctx := jwt.NewUserCtx(c.Request.Context(), claim.User)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
