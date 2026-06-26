package system

import (
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	"github.com/lihongsheng/go-admin/server/model/system"
	repoSys "github.com/lihongsheng/go-admin/server/repo/system"
	casbinUtil "github.com/lihongsheng/go-admin/server/utils/casbin"
)

// RoleService 角色业务接口
type RoleService interface {
	Create(req dtoSys.RoleCreateReq) (*system.SysRole, error)
	Update(req dtoSys.RoleUpdateReq) error
	Delete(id uint) error
	List() (*dtoSys.RoleListResp, error)
	Auth(req dtoSys.RoleAuthReq) error
	AuthDetail(id uint) (*dtoSys.RoleAuthDetailResp, error)
	SetDefaultRouter(roleID uint, defaultRouter string) error
}

// NewRoleService 构造 RoleService
func NewRoleService(
	roleRepo repoSys.RoleRepo,
	menuRepo repoSys.MenuRepo,
	apiRepo repoSys.ApiRepo,
	casbin casbinUtil.Port,
) RoleService {
	return &roleService{
		roleRepo: roleRepo,
		menuRepo: menuRepo,
		apiRepo:  apiRepo,
		casbin:   casbin,
	}
}

type roleService struct {
	roleRepo repoSys.RoleRepo
	menuRepo repoSys.MenuRepo
	apiRepo  repoSys.ApiRepo
	casbin   casbinUtil.Port
}

// DefaultRole 包级单例
var DefaultRole RoleService

func (s *roleService) Create(req dtoSys.RoleCreateReq) (*system.SysRole, error) {
	r := &system.SysRole{
		Name: req.Name, Remark: req.Remark,
		Status: req.Status, DefaultRouter: req.DefaultRouter,
	}
	if err := s.roleRepo.Create(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *roleService) Update(req dtoSys.RoleUpdateReq) error {
	patch := map[string]any{
		"name":           req.Name,
		"remark":         req.Remark,
		"status":         req.Status,
		"default_router": req.DefaultRouter,
	}
	return s.roleRepo.Update(req.ID, patch)
}

func (s *roleService) Delete(id uint) error {
	// 清理 Casbin 策略
	_ = s.casbin.RemoveRolePolicies(id)
	return s.roleRepo.Delete(id)
}

func (s *roleService) List() (*dtoSys.RoleListResp, error) {
	list, err := s.roleRepo.List()
	if err != nil {
		return nil, err
	}
	return &dtoSys.RoleListResp{List: list, Total: len(list)}, nil
}

// Auth 角色授权（菜单 + API）
// 在写入 sys_role_menus 前自动补全所有父级菜单 ID，
// 防止前端只回传半选叶子节点导致父菜单丢失（前端 v-tree 半选状态修复见 web 端）。
func (s *roleService) Auth(req dtoSys.RoleAuthReq) error {
	// 补全父级菜单 ID
	menuIDs, err := s.menuRepo.CompleteParentIDs(req.MenuIDs)
	if err != nil {
		return err
	}
	menus, err := s.menuRepo.FindByIDs(menuIDs)
	if err != nil {
		return err
	}
	apis, err := s.apiRepo.FindByIDs(req.ApiIDs)
	if err != nil {
		return err
	}
	if err := s.roleRepo.ReplaceMenus(req.RoleID, menus); err != nil {
		return err
	}
	// 同步 Casbin 策略（API 权限完全由 Casbin 管理，不再维护 sys_role_apis 表）
	items := make([][2]string, 0, len(apis))
	for _, a := range apis {
		items = append(items, [2]string{a.Path, a.Method})
	}
	return s.casbin.ReplaceRolePolicies(req.RoleID, items)
}

func (s *roleService) AuthDetail(id uint) (*dtoSys.RoleAuthDetailResp, error) {
	role, err := s.roleRepo.GetByIDWithMenus(id)
	if err != nil {
		return nil, err
	}
	menuIDs := make([]uint, 0, len(role.Menus))
	for _, m := range role.Menus {
		menuIDs = append(menuIDs, m.ID)
	}
	apiIDs := make([]uint, 0)
	for _, p := range s.casbin.GetRolePolicies(role.ID) {
		api, err := s.apiRepo.FindByPathMethod(p[0], p[1])
		if err == nil && api != nil {
			apiIDs = append(apiIDs, api.ID)
		}
	}
	return &dtoSys.RoleAuthDetailResp{MenuIDs: menuIDs, ApiIDs: apiIDs, DefaultRouter: role.DefaultRouter}, nil
}

func (s *roleService) SetDefaultRouter(roleID uint, defaultRouter string) error {
	return s.roleRepo.Update(roleID, map[string]any{"default_router": defaultRouter})
}
