package middleware

import (
	"strings"

	"go-admin/server/utils/jwt"
	"go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// CtxUserID / CtxUsername / CtxRoles 上下文 key
const (
	CtxUserID   = "uid"
	CtxUsername = "username"
	CtxRoles    = "roles"
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
		claim, err := jwt.Parse(raw)
		if err != nil {
			response.FailCode(c, response.CodeUnauthorized, "invalid token: "+err.Error())
			c.Abort()
			return
		}
		c.Set(CtxUserID, claim.UserID)
		c.Set(CtxUsername, claim.Username)
		c.Set(CtxRoles, claim.Roles)
		c.Next()
	}
}
