// Package client 封装 gRPC 客户端拨号逻辑。
//
// 学习项目用明文（insecure）传输；生产环境应换用 TLS credentials。
// token 非空时，会通过客户端拦截器在每次调用自动附加
// authorization: Bearer <token> metadata，配合服务端 Auth 拦截器。
package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Dial 连接到目标 gRPC 服务，返回可复用的 *grpc.ClientConn。
// grpc.NewClient 是新版推荐 API（替代已废弃的 grpc.Dial），默认懒连接。
func Dial(target, token string) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if token != "" {
		opts = append(opts, grpc.WithUnaryInterceptor(authUnary(token)))
		opts = append(opts, grpc.WithStreamInterceptor(authStream(token)))
	}
	return grpc.NewClient(target, opts...)
}

func authUnary(token string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func authStream(token string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string,
		streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
