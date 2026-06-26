package system

import (
	dtoSys "github.com/lihongsheng/go-admin/server/dto/system"
	"github.com/lihongsheng/go-admin/server/model/system"
	repoSys "github.com/lihongsheng/go-admin/server/repo/system"
)

// ApiService API 元数据业务接口
//
// TODO: ApiUpdate / ApiDelete 暂未联动 casbin_rule；
// 后续若要保证 API 路径变更后旧 policy 不残留，可在此处批量更新 casbin policy。
type ApiService interface {
	Create(req dtoSys.ApiCreateReq) (*system.SysApi, error)
	Update(req dtoSys.ApiUpdateReq) error
	Delete(id uint) error
	List(req dtoSys.ApiListReq) (*dtoSys.ApiListResp, error)
}

// NewApiService 构造 ApiService
func NewApiService(apiRepo repoSys.ApiRepo) ApiService {
	return &apiService{repo: apiRepo}
}

type apiService struct {
	repo repoSys.ApiRepo
}

// DefaultApi 包级单例
var DefaultApi ApiService

func (s *apiService) Create(req dtoSys.ApiCreateReq) (*system.SysApi, error) {
	a := &system.SysApi{Path: req.Path, Method: req.Method, Group: req.Group, Desc: req.Desc}
	if err := s.repo.Create(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *apiService) Update(req dtoSys.ApiUpdateReq) error {
	patch := map[string]any{
		"path":   req.Path,
		"method": req.Method,
		"group":  req.Group,
		"desc":   req.Desc,
	}
	return s.repo.Update(req.ID, patch)
}

func (s *apiService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *apiService) List(req dtoSys.ApiListReq) (*dtoSys.ApiListResp, error) {
	list, err := s.repo.List(req.Group)
	if err != nil {
		return nil, err
	}
	return &dtoSys.ApiListResp{List: list, Total: len(list)}, nil
}
