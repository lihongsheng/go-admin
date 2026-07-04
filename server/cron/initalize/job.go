package initalize

import (
  "context"
  "github.com/lihongsheng/go-admin/server/cron"
  log2 "github.com/lihongsheng/go-admin/server/log"
  "time"
)

func GetCronJobs() []cron.Job {
  return []cron.Job{
    {
      Spec:        "@every 1m",
      JobName:     "测试任务",
      Job:         TestJob,
      MaxExecTime: 2 * time.Minute,
    },
  }
}

func TestJob(ctx context.Context) error {
  log2.Info("TestJob", "time", time.Now())
  return nil
}
