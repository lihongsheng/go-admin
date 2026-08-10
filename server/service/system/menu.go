package system

import (
	dtoSys "go-admin/server/dto/system"
	"go-admin/server/model/system"
	repoSys "go-admin/server/repo/system"
)

// MenuService 菜单业务接口
type MenuService interface {
	Create(req dtoSys.MenuCreateReq) (*system.SysMenu, error)
	Update(req dtoSys.MenuUpdateReq) error
	// Delete 递归删除目标菜单及全部后代节点
	Delete(id uint) error
	// Tree 全量菜单树（管理员侧菜单管理）
	Tree() (*dtoSys.MenuTreeResp, error)
	// UserTree 指定用户的菜单树（运行时左侧菜单）
	UserTree(userID uint) ([]system.SysMenu, error)
}

// NewMenuService 构造 MenuService；userRepo 用于查用户对应角色
func NewMenuService(menuRepo repoSys.MenuRepo, userRepo repoSys.UserRepo) MenuService {
	return &menuService{menuRepo: menuRepo, userRepo: userRepo}
}

type menuService struct {
	menuRepo repoSys.MenuRepo
	userRepo repoSys.UserRepo
}

// DefaultMenu 包级单例
var DefaultMenu MenuService

func (s *menuService) Create(req dtoSys.MenuCreateReq) (*system.SysMenu, error) {
	t := req.Type
	if t == "" {
		t = system.MenuTypeMenu
	}
	m := &system.SysMenu{
		Type: t, ParentID: req.ParentID,
		Path: req.Path, Name: req.Name,
		Component: req.Component, Redirect: req.Redirect,
		Permission: req.Permission,
		Title:      req.Title, Icon: req.Icon,
		Sort: req.Sort, Hidden: req.Hidden, KeepAlive: req.KeepAlive,
		ApiRules: req.ApiRules,
	}
	if err := s.menuRepo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *menuService) Update(req dtoSys.MenuUpdateReq) error {
	patch := map[string]any{
		"type":       req.Type,
		"parent_id":  req.ParentID,
		"path":       req.Path,
		"name":       req.Name,
		"component":  req.Component,
		"redirect":   req.Redirect,
		"permission": req.Permission,
		"title":      req.Title,
		"icon":       req.Icon,
		"sort":       req.Sort,
		"hidden":     req.Hidden,
		"keep_alive": req.KeepAlive,
		"api_rules":  req.ApiRules,
	}
	return s.menuRepo.Update(req.ID, patch)
}

func (s *menuService) Delete(id uint) error {
	return s.menuRepo.DeleteCascade(id)
}

func (s *menuService) Tree() (*dtoSys.MenuTreeResp, error) {
	list, err := s.menuRepo.ListAll()
	if err != nil {
		return nil, err
	}
	return &dtoSys.MenuTreeResp{List: repoSys.BuildTree(list, 0)}, nil
}

// UserTree 取当前用户角色对应的菜单，组装为树
func (s *menuService) UserTree(userID uint) ([]system.SysMenu, error) {
	u, err := s.userRepo.GetByID(userID, true)
	if err != nil {
		return nil, err
	}
	roleIDs := make([]uint, 0, len(u.Roles))
	for _, r := range u.Roles {
		roleIDs = append(roleIDs, r.ID)
	}
	menus, err := s.menuRepo.MenusByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}
	return repoSys.BuildTree(menus, 0), nil
}
