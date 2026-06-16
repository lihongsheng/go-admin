package middleware

import (
	"go-admin/server/utils/casbin"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// CasbinAuth casbin RBAC 鉴权中间件，必须放在 JWTAuth 之后
func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, exists := c.Get(CtxUserID)
		if !exists {
			response.FailHTTP(c, 403, response.CodeForbidden, "missing user identity")
			c.Abort()
			return
		}
		uidUint, ok := uid.(uint)
		if !ok {
			response.FailHTTP(c, 403, response.CodeForbidden, "invalid user identity")
			c.Abort()
			return
		}
		// 用 FullPath（路由模板，如 /api/v1/system/user/:id）匹配策略
		path := c.FullPath()
		method := c.Request.Method
		allowed, err := casbin.Enforce(uidUint, path, method)
		if err != nil {
			response.FailHTTP(c, 500, response.CodeError, "authorization error: "+err.Error())
			c.Abort()
			return
		}
		if !allowed {
			response.FailHTTP(c, 403, response.CodeForbidden, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
