// product-service 启动入口。监听 :50052。
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"

	productv1 "learngrpc/gen/product/v1"
	"learngrpc/internal/interceptor"
	"learngrpc/internal/server"
	"learngrpc/services/product"
)

func main() {
	logger := log.New(os.Stdout, "[product-service] ", log.LstdFlags|log.Lmsgprefix)
	addr := server.EnvOr("PRODUCT_SERVICE_ADDR", ":50052")

	err := server.Run(context.Background(), server.Config{
		Name: "product-service",
		Addr: addr,
		Register: func(s *grpc.Server) {
			productv1.RegisterProductServiceServer(s, product.New())
		},
		ServerOpts: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(interceptor.UnaryRecovery(logger), interceptor.UnaryLogging(logger)),
			grpc.ChainStreamInterceptor(interceptor.StreamRecovery(logger), interceptor.StreamLogging(logger)),
		},
	})
	if err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
