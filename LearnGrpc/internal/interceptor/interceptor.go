// Package interceptor 提供一组可复用的服务端拦截器。
//
// 拦截器是 gRPC 的中间件机制：
//   - UnaryServerInterceptor  拦截普通一元调用
//   - StreamServerInterceptor 拦截流式调用
// 用 grpc.ChainUnaryInterceptor(...) 串联多个，顺序执行。
package interceptor

import (
	"context"
	"log"
	"runtime/debug"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// ===================== 日志 =====================

// UnaryLogging 记录每个一元调用的方法、状态码、耗时、来源地址。
func UnaryLogging(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		logger.Printf("unary  %-42s code=%-7s dur=%-8s from=%s",
			info.FullMethod, status.Code(err),
			time.Since(start).Round(time.Millisecond), clientAddr(ctx))
		return resp, err
	}
}

// StreamLogging 记录每个流式调用的方法与状态码。
func StreamLogging(logger *log.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		err := handler(srv, ss)
		logger.Printf("stream %-42s code=%-7s dur=%-8s from=%s",
			info.FullMethod, status.Code(err),
			time.Since(start).Round(time.Millisecond), clientAddr(ss.Context()))
		return err
	}
}

func clientAddr(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}
	return "unknown"
}

// ===================== panic 恢复 =====================

// UnaryRecovery 捕获 handler 中的 panic，避免单个请求把整个进程拖垮。
func UnaryRecovery(logger *log.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("panic in %s: %v\n%s", info.FullMethod, r, debug.Stack())
				err = status.Error(codes.Internal, "internal error")
			}
		}()
		return handler(ctx, req)
	}
}

// StreamRecovery 流式版本的 panic 恢复。
func StreamRecovery(logger *log.Logger) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.Printf("panic in %s: %v\n%s", info.FullMethod, r, debug.Stack())
				err = status.Error(codes.Internal, "internal error")
			}
		}()
		return handler(srv, ss)
	}
}

// ===================== 鉴权 =====================

// UnaryAuth 校验请求 metadata 里的 authorization: Bearer <token>。
// 仅当 validToken 非空时启用。
func UnaryAuth(validToken string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := authorize(ctx, validToken); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

// StreamAuth 流式版本的鉴权。
func StreamAuth(validToken string) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if err := authorize(ss.Context(), validToken); err != nil {
			return err
		}
		return handler(srv, ss)
	}
}

func authorize(ctx context.Context, validToken string) error {
	if validToken == "" {
		return nil // 未配置 token，放行
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "missing metadata")
	}
	values := md.Get("authorization")
	if len(values) == 0 {
		return status.Error(codes.Unauthenticated, "authorization token not provided")
	}
	token := strings.TrimPrefix(values[0], "Bearer ")
	if token != validToken {
		return status.Error(codes.Unauthenticated, "invalid token")
	}
	return nil
}
