// Package upload 文件上传工具包
//
// 支持本地存储、阿里云 OSS、腾讯云 COS。通过 config.yaml 中的 upload.drive 切换。
package upload

import (
	"context"
	"go-admin/server/config"
	"mime/multipart"
)

// OSS 对象存储接口
type OSS interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error)
	DeleteFile(ctx context.Context, key string) error
}

// NewOss 根据配置中 drive 字段返回对应的 OSS 实例
func NewOss(cfg config.Upload) OSS {
	switch cfg.Drive {
	case "aliyun-oss":
		return &AliyunOSS{Config: cfg.AliyunOSS}
	case "tencent-cos":
		return &TencentCOS{Config: cfg.TencentCOS}
	case "local":
		return &Local{Config: cfg.Local}
	default:
		return &Local{Config: cfg.Local}
	}
}
