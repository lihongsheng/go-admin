package base

import (
	"go-admin/server/global"
	"go-admin/server/utils/response"
	"go-admin/server/utils/upload"
	"strings"

	"github.com/gin-gonic/gin"
)

const maxUploadSize = 10 << 20 // 10 MB

// Upload 上传文件
// @Summary      上传文件（登录用户均可上传）
// @Tags         基础接口
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "上传文件（最大 10MB）"
// @Success      200   {object}  response.Body{data=object{url=string,name=string}}
// @Security     BearerAuth
// @Router       /api/v1/base/upload [post]
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "读取上传文件失败: "+err.Error())
		return
	}
	if file.Size > maxUploadSize {
		response.Fail(c, "文件大小不能超过 10MB")
		return
	}

	oss := upload.NewOss(global.Cfg.Upload)
	filename, err := oss.UploadFile(c.Request.Context(), file)
	if err != nil {
		response.Fail(c, "上传失败: "+err.Error())
		return
	}

	response.OK(c, gin.H{
		"url":  toURL(filename),
		"name": file.Filename,
	})
}

// toURL 将存储返回结果转为访问 URL：云存储返回完整 URL，本地存储返回文件名需拼接前缀
func toURL(raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	return "/uploads/" + raw
}
