package middleware

import (
	serviceBase "go-admin/server/service/base"
	"go-admin/server/utils/casbin"
	"go-admin/server/utils/jwt"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// CasbinAuth casbin 鉴权中间件，必须放在 JWTAuth 之后
// 从请求上下文获取用户信息，校验用户状态，逐一检查每个角色的 API 权限
func CasbinAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := jwt.GetUser(c.Request.Context())
		if err != nil {
			response.FailHTTP(c, 403, response.CodeForbidden, "missing user identity")
			c.Abort()
			return
		}

		// 通过 service 层查询用户并校验状态
		if _, err := serviceBase.Default.GetActiveUser(u.ID); err != nil {
			response.FailHTTP(c, 403, response.CodeForbidden, err.Error())
			c.Abort()
			return
		}

		if len(u.Role) == 0 {
			response.FailHTTP(c, 403, response.CodeForbidden, "forbidden")
			c.Abort()
			return
		}

		// 用 FullPath（路由模板，如 /api/v1/system/user/:id）匹配策略
		path := c.FullPath()
		method := c.Request.Method
		allowed := false
		for _, rid := range u.Role {
			ok, err := casbin.Enforce(uint(rid), path, method)
			if err != nil {
				response.FailHTTP(c, 500, response.CodeError, "authorization error: "+err.Error())
				c.Abort()
				return
			}
			if ok {
				allowed = true
				break
			}
		}
		if !allowed {
			response.FailHTTP(c, 403, response.CodeForbidden, "forbidden")
			c.Abort()
			return
		}
		c.Next()
	}
}
