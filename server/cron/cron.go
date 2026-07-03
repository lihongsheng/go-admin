package cron

import (
	"context"
	log2 "github.com/lihongsheng/go-admin/server/log"
	"github.com/pkg/errors"
	cron2 "github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"time"
)

type Job struct {
	Spec        string
	JobName     string
	Job         JobFun
	MaxExecTime time.Duration
}

type JobFun func(ctx context.Context) error

type Server struct {
	c   *cron2.Cron
	log *log
}

func NewCronServer() *Server {
	c := &Server{
		log: new(log),
		c:   cron2.New(cron2.WithSeconds(), cron2.WithLogger(new(log)), cron2.WithChain(cron2.SkipIfStillRunning(new(log)))),
	}
	_, _ = c.c.AddFunc("@every 5m", func() {
		log2.Info("NewCronStart", zap.Any("time", time.Now()))
	})
	c.c.Start()
	return c
}

func (c *Server) Start() error {
	return nil
}

func (c *Server) Stop() error {
	c.c.Stop()
	return nil
}

func (c *Server) AddJob(jobs ...Job) error {
	var err error
	for _, job := range jobs {
		err2 := c.addJobFun(job.Spec, job.JobName, job.Job, job.MaxExecTime)
		if err2 != nil {
			err = errors.Wrap(err2, "AddJobFun"+job.JobName+"Fail")
		}
	}
	return err
}

// AddJobFun
// 添加定时脚本
func (c *Server) addJobFun(spec string, jobName string, job JobFun, maxExecTime time.Duration) error {
	c.log.Info("AddJobFun", "jobName", jobName)
	_, err := c.c.AddFunc(spec, func() {
		ctx, cancel := context.WithTimeout(context.Background(), maxExecTime)
		defer cancel()
		defer func() {
			r := recover()
			if _, ok := r.(error); ok {
				c.log.Error(r.(error), jobName+"Fail")
			}
		}()
		c.log.Info(jobName, "date", time.Now().Format(time.DateTime))
		err := job(ctx)
		if err != nil {
			c.log.Error(err, jobName+"Fail")
		}
	})
	return err
}
