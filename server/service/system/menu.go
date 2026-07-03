package system

import (
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	"github.com/lihongsheng/go-admin/server/enum"
	"github.com/lihongsheng/go-admin/server/model/system"
	repoSys "github.com/lihongsheng/go-admin/server/repo/system"
)

// MenuService 菜单业务接口
type MenuService interface {
	Create(req dtoSys.MenuCreateReq) (*system.SysMenu, error)
	Update(req dtoSys.MenuUpdateReq) error
	// Delete 递归删除目标菜单及全部后代节点
	Delete(id uint) error
	// Tree 全量菜单树（管理员侧菜单管理）
	Tree() (*dtoSys.MenuTreeResp, error)
	// TreeBySystemType 按系统类型获取菜单树
	TreeBySystemType(systemType int) (*dtoSys.MenuTreeResp, error)
	// UserTree 指定用户的菜单树（运行时左侧菜单）
	// 非平台用户（商家/代理）按 systemType 返回全量菜单，不走角色过滤
	UserTree(userID uint, systemType enum.SystemType) ([]system.SysMenu, error)
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
		SystemType: req.SystemType,
		ApiRules:   req.ApiRules,
	}
	if err := s.menuRepo.Create(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *menuService) Update(req dtoSys.MenuUpdateReq) error {
	patch := map[string]any{
		"type":        req.Type,
		"parent_id":   req.ParentID,
		"path":        req.Path,
		"name":        req.Name,
		"component":   req.Component,
		"redirect":    req.Redirect,
		"permission":  req.Permission,
		"title":       req.Title,
		"icon":        req.Icon,
		"sort":        req.Sort,
		"hidden":      req.Hidden,
		"keep_alive":  req.KeepAlive,
		"system_type": req.SystemType,
		"api_rules":   req.ApiRules,
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

func (s *menuService) TreeBySystemType(systemType int) (*dtoSys.MenuTreeResp, error) {
	list, err := s.menuRepo.ListBySystemType(enum.SystemType(systemType))
	if err != nil {
		return nil, err
	}
	return &dtoSys.MenuTreeResp{List: repoSys.BuildTree(list, 0)}, nil
}

// UserTree 取当前用户角色对应的菜单，组装为树
// 非平台用户按 systemType 返回全量菜单，不走角色过滤
func (s *menuService) UserTree(userID uint, systemType enum.SystemType) ([]system.SysMenu, error) {
	var menus []system.SysMenu
	var err error

	// 平台用户：按角色加载菜单
	u, err2 := s.userRepo.GetByID(userID, true)
	if err2 != nil {
		return nil, err2
	}
	roleIDs := make([]uint, 0, len(u.Roles))
	for _, r := range u.Roles {
		roleIDs = append(roleIDs, r.ID)
	}
	menus, err = s.menuRepo.MenusByRoleIDs(roleIDs)

	if err != nil {
		return nil, err
	}
	return repoSys.BuildTree(menus, 0), nil
}
