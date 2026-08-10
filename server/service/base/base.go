// Package base base 模块业务服务
package base

import (
	"errors"

	dtoBase "go-admin/server/dto/base"
	"go-admin/server/enum"
	"go-admin/server/global"
	"go-admin/server/model/system"
	repoSys "go-admin/server/repo/system"
	"go-admin/server/utils/captcha"
	"go-admin/server/utils/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service base 模块对外服务接口
type Service interface {
	Captcha() (*dtoBase.CaptchaResp, error)
	Login(req dtoBase.LoginReq) (*dtoBase.LoginResp, error)
	Info(uid uint) (interface{}, error)
	GetActiveUser(uid uint) (*system.SysUser, error)
}

// NewService 构造 base.Service
func NewService(userRepo repoSys.UserRepo) Service {
	return &service{userRepo: userRepo}
}

type service struct {
	userRepo repoSys.UserRepo
}

// Default 包级单例，由 initialize 装配
var Default Service

func (s *service) Captcha() (*dtoBase.CaptchaResp, error) {
	id, b64, err := captcha.Generate()
	if err != nil {
		return nil, err
	}
	return &dtoBase.CaptchaResp{CaptchaID: id, CaptchaB64: b64}, nil
}

// 业务错误（暴露给 handler 直接作为 message 返回）
var (
	ErrCaptcha      = errors.New("captcha invalid or expired")
	ErrUserNotFound = errors.New("user not found")
	ErrUserDisabled = errors.New("user disabled")
	ErrWrongPwd     = errors.New("wrong password")
)

func (s *service) Login(req dtoBase.LoginReq) (*dtoBase.LoginResp, error) {
	if !captcha.Verify(req.CaptchaID, req.CaptchaCode) {
		return nil, ErrCaptcha
	}
	u, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if enum.UserStatus(u.Status) == enum.UserStatusDisabled {
		return nil, ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return nil, ErrWrongPwd
	}
	roleIDs := make([]int64, 0, len(u.Roles))
	for _, r := range u.Roles {
		roleIDs = append(roleIDs, int64(r.ID))
	}
	token, err := jwt.Sign(jwt.User{
		ID:       u.ID,
		Username: u.Username,
		Role:     roleIDs,
	}, global.Cfg.JWT)
	if err != nil {
		return nil, err
	}
	return &dtoBase.LoginResp{Token: token, User: *u}, nil
}

func (s *service) Info(uid uint) (interface{}, error) {
	return s.userRepo.GetByID(uid, true)
}

// GetActiveUser 查询用户并检查是否被禁用
func (s *service) GetActiveUser(uid uint) (*system.SysUser, error) {
	u, err := s.userRepo.GetByID(uid, false)
	if err != nil {
		return nil, err
	}
	if enum.UserStatus(u.Status) == enum.UserStatusDisabled {
		return nil, ErrUserDisabled
	}
	return u, nil
}
