// Package install 安装向导 HTTP 接口
//
// 业务流程在 service/install 包；handler 只做参数绑定、SSE 适配。
package install

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/lihongsheng/go-admin/server/core/installer"
	dtoInstall "github.com/lihongsheng/go-admin/server/dto/install"
	"github.com/lihongsheng/go-admin/server/global"
	serviceInstall "github.com/lihongsheng/go-admin/server/service/install"
	"github.com/lihongsheng/go-admin/server/utils/response"

	"github.com/gin-gonic/gin"
)

// Status GET /install/status
func Status(c *gin.Context) {
	response.OK(c, serviceInstall.Default.Status())
}

// CheckDB POST /install/check-db
func CheckDB(c *gin.Context) {
	var in dtoInstall.CheckDBReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, "invalid db config: "+err.Error())
		return
	}
	response.OK(c, serviceInstall.Default.CheckDB(in))
}

// Init POST /install/init —— 同步执行；返回最终结果
func Init(c *gin.Context) {
	var req dtoInstall.InitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid body: "+err.Error())
		return
	}
	if global.Installed.Load() {
		response.FailCode(c, response.CodeAlreadyDone, "already installed")
		return
	}

	// 收集步骤进度（同步模式下安装完一次性返回）
	steps := make([]installer.Step, 0, 32)
	ch := make(chan installer.Step, 32)
	done := make(chan struct{})
	go func() {
		for s := range ch {
			steps = append(steps, s)
		}
		close(done)
	}()

	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Minute)
	defer cancel()

	createdDB, err := serviceInstall.Default.Install(ctx, req, ch)
	close(ch)
	<-done

	if err != nil {
		response.OK(c, gin.H{"ok": false, "error": err.Error(), "steps": steps})
		return
	}
	if err := serviceInstall.Default.AfterInstalled(req.DB); err != nil {
		response.OK(c, gin.H{"ok": false, "error": err.Error(), "steps": steps})
		return
	}
	response.OK(c, gin.H{"ok": true, "created_db": createdDB, "steps": steps})
}

// Stream POST /install/stream —— SSE 流式安装；与 Init 二选一
func Stream(c *gin.Context) {
	var req dtoInstall.InitReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "invalid body: "+err.Error())
		return
	}
	if global.Installed.Load() {
		response.FailCode(c, response.CodeAlreadyDone, "already installed")
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	send := func(event string, payload interface{}) {
		b, _ := json.Marshal(payload)
		_, _ = io.WriteString(c.Writer, fmt.Sprintf("event: %s\ndata: %s\n\n", event, b))
		c.Writer.Flush()
	}

	ch := make(chan installer.Step, 32)
	errCh := make(chan error, 1)
	var createdDB bool

	go func() {
		var err error
		createdDB, err = serviceInstall.Default.Install(c.Request.Context(), req, ch)
		errCh <- err
		close(ch)
	}()

	for s := range ch {
		send("step", s)
	}
	if err := <-errCh; err != nil {
		send("error", gin.H{"error": err.Error()})
		return
	}
	// install 成功后再把"是否新建库"事件回放一次（保持与重构前的事件名兼容）
	send("ensure_db", gin.H{"created": createdDB})

	if err := serviceInstall.Default.AfterInstalled(req.DB); err != nil {
		send("error", gin.H{"error": err.Error()})
		return
	}
	send("done", gin.H{"ok": true, "created_db": createdDB})
}
