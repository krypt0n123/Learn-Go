// user-service 启动入口。监听 :50051。
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"

	userv1 "learngrpc/gen/user/v1"
	"learngrpc/internal/interceptor"
	"learngrpc/internal/server"
	"learngrpc/services/user"
)

func main() {
	logger := log.New(os.Stdout, "[user-service] ", log.LstdFlags|log.Lmsgprefix)
	addr := server.EnvOr("USER_SERVICE_ADDR", ":50051")

	err := server.Run(context.Background(), server.Config{
		Name: "user-service",
		Addr: addr,
		Register: func(s *grpc.Server) {
			userv1.RegisterUserServiceServer(s, user.New())
		},
		// 拦截器按顺序执行：先 recovery 兜底，再 logging 记录。
		ServerOpts: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(interceptor.UnaryRecovery(logger), interceptor.UnaryLogging(logger)),
			grpc.ChainStreamInterceptor(interceptor.StreamRecovery(logger), interceptor.StreamLogging(logger)),
		},
	})
	if err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
