package system

import "github.com/lihongsheng/go-admin/server/model/system"

// ApiCreateReq 新增 API
type ApiCreateReq struct {
	Path   string `json:"path"   binding:"required"`
	Method string `json:"method" binding:"required"`
	Group  string `json:"group"`
	Desc   string `json:"desc"`
}

// ApiUpdateReq 更新 API
type ApiUpdateReq struct {
	ID     uint   `json:"id"     binding:"required"`
	Path   string `json:"path"   binding:"required"`
	Method string `json:"method" binding:"required"`
	Group  string `json:"group"`
	Desc   string `json:"desc"`
}

// ApiListReq API 列表过滤
type ApiListReq struct {
	Group string `form:"group"`
}

// ApiListResp API 列表响应
type ApiListResp struct {
	List  []system.SysApi `json:"list"`
	Total int             `json:"total"`
}
