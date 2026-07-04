package server

import (
	"context"
	"errors"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type App struct {
	ctx         context.Context
	cancel      func()
	instance    []Server
	stopTimeOut time.Duration
}

func NewApp(stopTimeOut time.Duration, instance ...Server) *App {
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		ctx:         ctx,
		cancel:      cancel,
		instance:    instance,
		stopTimeOut: stopTimeOut,
	}
}

func (a *App) Run() error {
	// errgroup 来控制服务启动，切字段 panic兜底
	// 注册ctx 当有一个服务启动失败的告知其他服务停止
	eg, ctx := errgroup.WithContext(a.ctx)
	// 等待所有服务启动完毕
	wg := sync.WaitGroup{}
	for _, instance := range a.instance {
		srv := instance
		eg.Go(func() error {
			<-ctx.Done() // 通过NewApp创建的ctx，等待服务停止信号
			stopCtx, cancel := context.WithTimeout(ctx, a.stopTimeOut)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done() //
			return srv.Start(ctx)
		})
	}
	// 等待所有服务启动完毕
	wg.Wait()
	// 注册退出信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	eg.Go(func() error {
		select {
		case <-ctx.Done(): // 如果服务启动失败，则直接退出
			return nil
		case <-c: // 通过信号量 来感知退出，比如Linux 的kill -9 进程
			return a.Stop()
		}
	})
	// 等待所有服务退出
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	// 通过NewApp创建的ctx 通知服务停止
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
