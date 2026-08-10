// Package base base 模块 DTO：登录 / 验证码 / 当前用户
package base

import "go-admin/server/model/system"

// LoginReq 登录请求
type LoginReq struct {
	Username    string `json:"username"     binding:"required"`
	Password    string `json:"password"     binding:"required"`
	CaptchaID   string `json:"captcha_id"   binding:"required"`
	CaptchaCode string `json:"captcha_code" binding:"required"`
}

// LoginResp 登录成功响应
type LoginResp struct {
	Token string         `json:"token"`
	User  system.SysUser `json:"user"`
}

// CaptchaResp 验证码响应
type CaptchaResp struct {
	CaptchaID  string `json:"captcha_id"`
	CaptchaB64 string `json:"captcha_b64"`
}
