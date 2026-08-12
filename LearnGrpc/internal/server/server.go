// Package server 封装 gRPC 服务端启动 + 优雅关闭的通用逻辑，
// 三个服务复用同一套骨架。
package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

// Config 描述一个 gRPC 服务的启动配置。
type Config struct {
	Name            string                   // 服务名，仅用于日志
	Addr            string                   // 监听地址，如 ":50051"
	Register        func(*grpc.Server)       // 把具体实现注册进 server
	ServerOpts      []grpc.ServerOption      // 透传给 grpc.NewServer（一般用来挂拦截器）
	ShutdownTimeout time.Duration            // 优雅关闭最长等待时间，0=默认 10s
}

// Run 启动 gRPC 服务，并在收到 SIGINT/SIGTERM 时优雅关闭。
//
// 优雅关闭（GracefulStop）会等待已开始的请求处理完再退出，
// 相比直接 Stop() 能避免打断了正在处理的调用。
func Run(ctx context.Context, cfg Config) error {
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 10 * time.Second
	}

	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", cfg.Addr, err)
	}

	srv := grpc.NewServer(cfg.ServerOpts...)
	cfg.Register(srv)

	errCh := make(chan error, 1)
	go func() {
		log.Printf("[%s] gRPC server listening on %s", cfg.Name, lis.Addr())
		if err := srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			errCh <- err
		}
	}()

	// 监听中断信号
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Printf("[%s] received shutdown signal, gracefully stopping...", cfg.Name)
		stopped := make(chan struct{})
		go func() {
			srv.GracefulStop()
			close(stopped)
		}()
		select {
		case <-stopped:
			log.Printf("[%s] stopped cleanly", cfg.Name)
		case <-time.After(cfg.ShutdownTimeout):
			log.Printf("[%s] graceful stop timed out, forcing stop", cfg.Name)
			srv.Stop()
		}
		return nil
	}
}

// EnvOr 从环境变量读取，没有则返回默认值。
func EnvOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
