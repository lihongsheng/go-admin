package system

import (
	"sort"

	"github.com/lihongsheng/go-admin/server/model/system"

	"gorm.io/gorm"
)

// MenuRepo 菜单仓储接口
type MenuRepo interface {
	Create(m *system.SysMenu) error
	Update(id uint, patch map[string]any) error
	// DeleteCascade 递归删除目标节点及其所有子孙节点
	DeleteCascade(id uint) error
	FindByIDs(ids []uint) ([]system.SysMenu, error)
	// ListAll 查全量菜单
	ListAll() ([]system.SysMenu, error)
	// CompleteParentIDs 给定一组菜单 ID，递归补全其所有父级 ID 后返回去重列表
	CompleteParentIDs(ids []uint) ([]uint, error)
	// MenusByRoleIDs 通过角色 ID 列表查出全部关联菜单（已按 sort 排序、按 ID 去重）
	MenusByRoleIDs(roleIDs []uint) ([]system.SysMenu, error)
}

// NewMenuRepo 构造 MenuRepo
func NewMenuRepo(db *gorm.DB) MenuRepo { return &menuRepo{db: db} }

type menuRepo struct{ db *gorm.DB }

func (r *menuRepo) Create(m *system.SysMenu) error {
	return r.db.Create(m).Error
}

func (r *menuRepo) Update(id uint, patch map[string]any) error {
	return r.db.Model(&system.SysMenu{}).Where("id=?", id).Updates(patch).Error
}

// DeleteCascade 递归收集所有后代节点 ID 后一次性删除
func (r *menuRepo) DeleteCascade(id uint) error {
	var allMenus []system.SysMenu
	if err := r.db.Select("id", "parent_id").Find(&allMenus).Error; err != nil {
		return err
	}
	// children 索引
	childMap := make(map[uint][]uint, len(allMenus))
	for _, m := range allMenus {
		childMap[m.ParentID] = append(childMap[m.ParentID], m.ID)
	}
	// BFS 收集自身 + 所有后代
	delIDs := []uint{id}
	queue := []uint{id}
	for len(queue) > 0 {
		head := queue[0]
		queue = queue[1:]
		for _, child := range childMap[head] {
			delIDs = append(delIDs, child)
			queue = append(queue, child)
		}
	}
	return r.db.Delete(&system.SysMenu{}, delIDs).Error
}

func (r *menuRepo) FindByIDs(ids []uint) ([]system.SysMenu, error) {
	var menus []system.SysMenu
	if len(ids) == 0 {
		return menus, nil
	}
	if err := r.db.Where("id IN ?", ids).Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

func (r *menuRepo) ListAll() ([]system.SysMenu, error) {
	var list []system.SysMenu
	if err := r.db.Order("sort").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *menuRepo) CompleteParentIDs(ids []uint) ([]uint, error) {
	if len(ids) == 0 {
		return ids, nil
	}
	var allMenus []system.SysMenu
	if err := r.db.Select("id", "parent_id").Find(&allMenus).Error; err != nil {
		return nil, err
	}
	parentOf := make(map[uint]uint, len(allMenus))
	for _, m := range allMenus {
		parentOf[m.ID] = m.ParentID
	}
	set := make(map[uint]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	for _, id := range ids {
		cur := id
		for {
			parent, ok := parentOf[cur]
			if !ok || parent == 0 {
				break
			}
			set[parent] = struct{}{}
			cur = parent
		}
	}
	out := make([]uint, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out, nil
}

func (r *menuRepo) MenusByRoleIDs(roleIDs []uint) ([]system.SysMenu, error) {
	if len(roleIDs) == 0 {
		return []system.SysMenu{}, nil
	}
	var menuIDs []uint
	if err := r.db.Table("sys_role_menus").
		Where("sys_role_id IN ?", roleIDs).
		Pluck("sys_menu_id", &menuIDs).Error; err != nil {
		return nil, err
	}
	if len(menuIDs) == 0 {
		return []system.SysMenu{}, nil
	}
	var menus []system.SysMenu
	if err := r.db.Where("id IN ?", menuIDs).Order("sort").Find(&menus).Error; err != nil {
		return nil, err
	}
	// 按 ID 去重（同 ID 保留 sort 更小的那条，理论上不重复）
	seen := map[uint]system.SysMenu{}
	for _, m := range menus {
		if exist, ok := seen[m.ID]; !ok || m.Sort < exist.Sort {
			seen[m.ID] = m
		}
	}
	out := make([]system.SysMenu, 0, len(seen))
	for _, m := range seen {
		out = append(out, m)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Sort < out[j].Sort })
	return out, nil
}

// BuildTree 将扁平菜单列表按 parent_id 组装为树。
// 通用工具，给 service 层在不同场景下复用（用户菜单树 / 全量菜单树）。
func BuildTree(list []system.SysMenu, parent uint) []system.SysMenu {
	out := []system.SysMenu{}
	for _, m := range list {
		if m.ParentID == parent {
			m.Children = BuildTree(list, m.ID)
			out = append(out, m)
		}
	}
	return out
}
