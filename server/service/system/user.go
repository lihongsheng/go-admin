// Package system system 模块业务服务
package system

import (
	dtoSys "go-admin/server/dto/system"
	"go-admin/server/model/system"
	repoSys "go-admin/server/repo/system"
	casbinUtil "go-admin/server/utils/casbin"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户业务接口
type UserService interface {
	Create(req dtoSys.UserCreateReq) (*system.SysUser, error)
	Update(req dtoSys.UserUpdateReq) error
	Delete(id uint) error
	List(req dtoSys.UserListReq) (*dtoSys.UserListResp, error)
}

// NewUserService 构造 UserService
func NewUserService(userRepo repoSys.UserRepo, casbin casbinUtil.Port) UserService {
	return &userService{repo: userRepo, casbin: casbin}
}

type userService struct {
	repo   repoSys.UserRepo
	casbin casbinUtil.Port
}

// DefaultUser 包级单例
var DefaultUser UserService

func (s *userService) Create(req dtoSys.UserCreateReq) (*system.SysUser, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := &system.SysUser{
		Username: req.Username, Password: string(hash),
		Nickname: req.Nickname, Email: req.Email,
		Phone: req.Phone, Status: req.Status,
	}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	if len(req.RoleIDs) > 0 {
		roles, err := s.repo.FindRolesByIDs(req.RoleIDs)
		if err != nil {
			return nil, err
		}
		if err := s.repo.ReplaceRoles(u.ID, roles); err != nil {
			return nil, err
		}
		if err := s.syncCasbinRoles(u.ID); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (s *userService) Update(req dtoSys.UserUpdateReq) error {
	patch := map[string]any{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"status":   req.Status,
	}
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		patch["password"] = string(hash)
	}
	if err := s.repo.Update(req.ID, patch); err != nil {
		return err
	}
	if req.RoleIDs != nil {
		roles, err := s.repo.FindRolesByIDs(req.RoleIDs)
		if err != nil {
			return err
		}
		if err := s.repo.ReplaceRoles(req.ID, roles); err != nil {
			return err
		}
		if err := s.syncCasbinRoles(req.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *userService) Delete(id uint) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	// 顺手清理 Casbin g 策略，避免 user-role 绑定脏数据
	return s.casbin.RemoveUserRoles(id)
}

func (s *userService) List(req dtoSys.UserListReq) (*dtoSys.UserListResp, error) {
	list, total, err := s.repo.List(req.Keyword, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &dtoSys.UserListResp{List: list, Total: total}, nil
}

// syncCasbinRoles 拉最新的用户角色码同步到 casbin g 策略
func (s *userService) syncCasbinRoles(uid uint) error {
	u, err := s.repo.GetByID(uid, true)
	if err != nil {
		return err
	}
	codes := make([]string, 0, len(u.Roles))
	for _, r := range u.Roles {
		codes = append(codes, r.Code)
	}
	return s.casbin.ReplaceUserRoles(uid, codes)
}
