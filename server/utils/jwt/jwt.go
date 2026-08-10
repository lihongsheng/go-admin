// Package jwt JWT 工具
package jwt

import (
	"context"
	"errors"
	jwtv5 "github.com/golang-jwt/jwt/v5"
	"go-admin/server/utils/jwt/config"
	"time"
)

type Claims struct {
	jwtv5.RegisteredClaims
	User
}
type User struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Role     []int64 `json:"role"`
}

type JwtContextUserKey struct{}

func GetUser(ctx context.Context) (User, error) {
	u := ctx.Value(JwtContextUserKey{})
	if u == nil {
		return User{}, errors.New("no user in context")
	}
	return u.(User), nil
}

func NewUserCtx(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, JwtContextUserKey{}, user)
}

// Sign 生成 token
func Sign(user User, jwtConfig config.JWT) (string, error) {
	cfg := jwtConfig
	c := Claims{
		User: user,
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
func Parse(tokenStr string, jwtConfig config.JWT) (*Claims, error) {
	cfg := jwtConfig
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
