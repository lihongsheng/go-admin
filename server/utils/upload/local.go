package upload

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/lihongsheng/go-admin/server/config"
	"github.com/lihongsheng/go-admin/server/log"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var mu sync.Mutex

type Local struct {
	Config config.Local
}

// UploadFile 保存文件到本地，返回文件名
func (l *Local) UploadFile(ctx context.Context, file *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(file.Filename)
	name := strings.TrimSuffix(file.Filename, ext)

	hash := md5.Sum([]byte(name))
	name = hex.EncodeToString(hash[:])[:8]
	filename := name + "_" + time.Now().Format("20060102150405") + ext

	if err := os.MkdirAll(l.Config.Path, os.ModePerm); err != nil {
		log.FromContext(ctx).Error("function os.MkdirAll() failed", "err", err.Error())
		return "", errors.New("function os.MkdirAll() failed, err:" + err.Error())
	}

	p := filepath.Join(l.Config.Path, filename)

	f, err := file.Open()
	if err != nil {
		log.FromContext(ctx).Error("function file.Open() failed", "err", err.Error())
		return "", errors.New("function file.Open() failed, err:" + err.Error())
	}
	defer f.Close()

	out, err := os.Create(p)
	if err != nil {
		log.FromContext(ctx).Error("function os.Create() failed", "err", err.Error())
		return "", errors.New("function os.Create() failed, err:" + err.Error())
	}
	defer out.Close()

	if _, err := io.Copy(out, f); err != nil {
		log.FromContext(ctx).Error("function io.Copy() failed", "err", err.Error())
		return "", errors.New("function io.Copy() failed, err:" + err.Error())
	}

	return filename, nil
}

// DeleteFile 删除本地文件
func (l *Local) DeleteFile(ctx context.Context, key string) error {
	if key == "" {
		return errors.New("key不能为空")
	}
	if strings.Contains(key, "..") || strings.ContainsAny(key, `\/:*?"<>|`) {
		return errors.New("非法的key")
	}

	p := filepath.Join(l.Config.Path, filepath.Base(key))
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return errors.New("文件不存在")
	}

	mu.Lock()
	defer mu.Unlock()
	return os.Remove(p)
}
