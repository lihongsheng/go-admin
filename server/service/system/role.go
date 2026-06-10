package system

import (
	dtoSys "go-admin/server/dto/system"
	"go-admin/server/model/system"
	repoSys "go-admin/server/repo/system"
	casbinUtil "go-admin/server/utils/casbin"
)

// RoleService 角色业务接口
type RoleService interface {
	Create(req dtoSys.RoleCreateReq) (*system.SysRole, error)
	Update(req dtoSys.RoleUpdateReq) error
	Delete(id uint) error
	List() (*dtoSys.RoleListResp, error)
	Auth(req dtoSys.RoleAuthReq) error
	AuthDetail(id uint) (*dtoSys.RoleAuthDetailResp, error)
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
		Name: req.Name, Code: req.Code,
		Remark: req.Remark, Status: req.Status,
	}
	if err := s.roleRepo.Create(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *roleService) Update(req dtoSys.RoleUpdateReq) error {
	// 检测 code 变更：需要迁移 Casbin 策略
	old, err := s.roleRepo.GetByID(req.ID)
	if err == nil && old.Code != req.Code {
		// 旧角色的 p 策略需移除（新 code 还没策略）；同时迁移 g 关联
		_ = s.casbin.RemoveRolePolicies(old.Code)
		_ = s.casbin.MigrateRoleCode(old.Code, req.Code)
	}
	patch := map[string]any{
		"name":   req.Name,
		"code":   req.Code,
		"remark": req.Remark,
		"status": req.Status,
	}
	return s.roleRepo.Update(req.ID, patch)
}

func (s *roleService) Delete(id uint) error {
	role, err := s.roleRepo.GetByID(id)
	if err == nil {
		_ = s.casbin.RemoveRolePolicies(role.Code)
		_ = s.casbin.RemoveRoleFromUsers(role.Code)
	}
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
	role, err := s.roleRepo.GetByID(req.RoleID)
	if err != nil {
		return err
	}
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
	if err := s.roleRepo.ReplaceMenus(role.ID, menus); err != nil {
		return err
	}
	if err := s.roleRepo.ReplaceApis(role.ID, apis); err != nil {
		return err
	}
	// 同步 Casbin 策略
	items := make([][2]string, 0, len(apis))
	for _, a := range apis {
		items = append(items, [2]string{a.Path, a.Method})
	}
	return s.casbin.ReplaceRolePolicies(role.Code, items)
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
	for _, p := range s.casbin.GetRolePolicies(role.Code) {
		api, err := s.apiRepo.FindByPathMethod(p[0], p[1])
		if err == nil && api != nil {
			apiIDs = append(apiIDs, api.ID)
		}
	}
	return &dtoSys.RoleAuthDetailResp{MenuIDs: menuIDs, ApiIDs: apiIDs}, nil
}
