package upload

import (
	"context"
	"errors"
	"fmt"
	"github.com/lihongsheng/go-admin/server/config"
	"github.com/lihongsheng/go-admin/server/log"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type TencentCOS struct {
	Config config.TencentCOS
}

// UploadFile 上传文件到腾讯云 COS，返回访问 URL
func (t *TencentCOS) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	client := newCosClient(t.Config)

	f, err := file.Open()
	if err != nil {
		log.FromContext(ctx).Error("file.Open failed", "err", err.Error())
		return "", errors.New("file.Open failed: " + err.Error())
	}
	defer f.Close()

	fileKey := fmt.Sprintf("%s/%d_%s", t.Config.PathPrefix, time.Now().Unix(), file.Filename)

	if _, err := client.Object.Put(ctx, fileKey, f, nil); err != nil {
		log.FromContext(ctx).Error("cos.Object.Put failed", "err", err.Error())
		return "", errors.New("cos.Object.Put failed: " + err.Error())
	}

	return t.Config.BaseURL + "/" + fileKey, nil
}

// DeleteFile 从腾讯云 COS 删除文件
func (t *TencentCOS) DeleteFile(ctx context.Context, key string) error {
	client := newCosClient(t.Config)
	name := t.Config.PathPrefix + "/" + key
	if _, err := client.Object.Delete(ctx, name); err != nil {
		log.FromContext(ctx).Error("cos.Object.Delete failed", "err", err.Error())
		return errors.New("cos.Object.Delete failed: " + err.Error())
	}
	return nil
}

func newCosClient(cfg config.TencentCOS) *cos.Client {
	urlStr, _ := url.Parse("https://" + cfg.Bucket + ".cos." + cfg.Region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	return cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretID,
			SecretKey: cfg.SecretKey,
		},
	})
}
