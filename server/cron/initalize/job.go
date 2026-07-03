package initalize

import (
	"context"
	"github.com/lihongsheng/go-admin/server/cron"
	log2 "github.com/lihongsheng/go-admin/server/log"
	"time"
)

func getCron() []cron.Job {
	return []cron.Job{
		{
			Spec:        "@every 1m",
			JobName:     "生成设备ota升级任务",
			Job:         TestJob,
			MaxExecTime: 2 * time.Minute,
		},
	}
}

func TestJob(ctx context.Context) error {
	log2.Info("TestJob", "time", time.Now())
	return nil
}

func Init() *cron.Server {
	jobs := getCron()
	c := cron.NewCronServer()
	for _, job := range jobs {
		err := c.AddJob(job)
		if err != nil {
			log2.Error("Init", "err", err)
		}
	}
	go func() {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}()
	return c
}
