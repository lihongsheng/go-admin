package upload

import (
	"context"
	"errors"
	"github.com/lihongsheng/go-admin/server/config"
	"github.com/lihongsheng/go-admin/server/log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliyunOSS struct {
	Config config.AliyunOSS
}

func (a *AliyunOSS) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	bucket, err := newAliyunBucket(a.Config)
	if err != nil {
		log.FromContext(ctx).Error("AliyunOSS.NewBucket failed", "err", err.Error())
		return "", errors.New("AliyunOSS.NewBucket failed: " + err.Error())
	}

	f, err := file.Open()
	if err != nil {
		log.FromContext(ctx).Error("file.Open failed", "err", err.Error())
		return "", errors.New("file.Open failed: " + err.Error())
	}
	defer f.Close()

	// 路径格式：{BasePath}/uploads/{yyyy-mm-dd}/{filename}
	ext := filepath.Ext(file.Filename)
	key := a.Config.BasePath + "/uploads/" + time.Now().Format("2006-01-02") + "/" + time.Now().Format("150405") + "_" + ext

	if err := bucket.PutObject(key, f); err != nil {
		log.FromContext(ctx).Error("bucket.PutObject failed", "err", err.Error())
		return "", errors.New("bucket.PutObject failed: " + err.Error())
	}

	return a.Config.BucketUrl + "/" + key, nil
}

func (a *AliyunOSS) DeleteFile(ctx context.Context, key string) error {
	bucket, err := newAliyunBucket(a.Config)
	if err != nil {
		log.FromContext(ctx).Error("AliyunOSS.NewBucket failed", "err", err.Error())
		return errors.New("AliyunOSS.NewBucket failed: " + err.Error())
	}

	if err := bucket.DeleteObject(key); err != nil {
		log.FromContext(ctx).Error("bucket.DeleteObject failed", "err", err.Error())
		return errors.New("bucket.DeleteObject failed: " + err.Error())
	}
	return nil
}

func newAliyunBucket(cfg config.AliyunOSS) (*oss.Bucket, error) {
	client, err := oss.New(cfg.Endpoint, cfg.AccessKey, cfg.AccessSecret)
	if err != nil {
		return nil, err
	}
	return client.Bucket(cfg.BucketName)
}
