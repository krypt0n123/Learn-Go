// order-service 启动入口。监听 :50053。
//
// 它对外是服务端，对下游又是客户端：启动时分别拨号连接
// user-service(:50051) 和 product-service(:50052)。
package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"

	orderv1 "learngrpc/gen/order/v1"
	productv1 "learngrpc/gen/product/v1"
	userv1 "learngrpc/gen/user/v1"
	"learngrpc/internal/client"
	"learngrpc/internal/interceptor"
	"learngrpc/internal/server"
	"learngrpc/services/order"
)

func main() {
	logger := log.New(os.Stdout, "[order-service] ", log.LstdFlags|log.Lmsgprefix)
	addr := server.EnvOr("ORDER_SERVICE_ADDR", ":50053")
	userAddr := server.EnvOr("USER_SERVICE_ADDR", "localhost:50051")
	productAddr := server.EnvOr("PRODUCT_SERVICE_ADDR", "localhost:50052")

	// 作为客户端连接下游服务。内部服务间调用不携带鉴权 token。
	userConn, err := client.Dial(userAddr, "")
	if err != nil {
		logger.Fatalf("dial user-service: %v", err)
	}
	defer userConn.Close()
	productConn, err := client.Dial(productAddr, "")
	if err != nil {
		logger.Fatalf("dial product-service: %v", err)
	}
	defer productConn.Close()

	clients := order.Clients{
		User:    userv1.NewUserServiceClient(userConn),
		Product: productv1.NewProductServiceClient(productConn),
	}

	// 对外开放的鉴权 token（演示用；生产环境应换 TLS + 真正的鉴权）。
	authToken := server.EnvOr("AUTH_TOKEN", "dev-secret")

	err = server.Run(context.Background(), server.Config{
		Name: "order-service",
		Addr: addr,
		Register: func(s *grpc.Server) {
			orderv1.RegisterOrderServiceServer(s, order.New(clients))
		},
		// 拦截器顺序：recovery 兜底 -> logging 记录(含被拒绝的调用) -> auth 鉴权 -> 业务
		ServerOpts: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(
				interceptor.UnaryRecovery(logger),
				interceptor.UnaryLogging(logger),
				interceptor.UnaryAuth(authToken),
			),
			grpc.ChainStreamInterceptor(
				interceptor.StreamRecovery(logger),
				interceptor.StreamLogging(logger),
				interceptor.StreamAuth(authToken),
			),
		},
	})
	if err != nil {
		logger.Fatalf("server error: %v", err)
	}
}
