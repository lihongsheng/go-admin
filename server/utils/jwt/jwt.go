// Package jwt JWT 工具
package jwt

import (
	"errors"
	"time"

	"go-admin/server/global"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	Roles    []string `json:"roles"`
	jwtv5.RegisteredClaims
}

// Sign 生成 token
func Sign(uid uint, username string, roles []string) (string, error) {
	cfg := global.Cfg.JWT
	c := Claims{
		UserID: uid, Username: username, Roles: roles,
		RegisteredClaims: jwtv5.RegisteredClaims{
			Issuer:    cfg.Issuer,
			IssuedAt:  jwtv5.NewNumericDate(time.Now()),
			ExpiresAt: jwtv5.NewNumericDate(time.Now().Add(time.Duration(cfg.Expire) * time.Second)),
		},
	}
	t := jwtv5.NewWithClaims(jwtv5.SigningMethodHS256, c)
	return t.SignedString([]byte(cfg.Secret))
}

// Parse 解析 token
func Parse(tokenStr string) (*Claims, error) {
	cfg := global.Cfg.JWT
	t, err := jwtv5.ParseWithClaims(tokenStr, &Claims{}, func(t *jwtv5.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	c, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return nil, errors.New("invalid token")
	}
	return c, nil
}
