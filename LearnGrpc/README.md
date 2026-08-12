# LearnGrpc —— 用一个微服务示例搞懂 gRPC 全流程

一个最小可运行的 gRPC 微服务学习项目，包含三个服务，演示：

- **四种 RPC 调用类型**：一元（Unary）、服务端流、客户端流、双向流
- **服务间调用**：`order-service` 在处理请求时，作为 gRPC 客户端去调用 `user-service` 与 `product-service`
- **官方推荐的最佳实践**：拦截器、状态码错误、超时/截止时间传播、优雅关闭、连接复用

> 📖 详细的一步步讲解见 [docs/guide.md](docs/guide.md)

## 架构

```
                ┌──────────────┐
                │  demo-client │  (gRPC 客户端)
                └──────┬───────┘
                       │ CreateOrder / GetOrder / StreamOrderUpdates
                       ▼
                ┌──────────────┐
                │ order-service│ :50053  (服务端 + 对下游是客户端)
                └───┬──────┬───┘
        GetUser │      │ CheckStock / GetProduct
                    ▼        ▼
        ┌──────────────┐  ┌─────────────────┐
        │ user-service │  │ product-service │
        │    :50051    │  │     :50052      │
        └──────────────┘  └─────────────────┘
```

## 目录结构

```
.
├── proto/            # .proto 接口定义（IDL）
├── gen/              # protoc 生成的 Go 代码（不要手改）
├── internal/
│   ├── server/       # 服务端启动 + 优雅关闭通用骨架
│   ├── interceptor/  # 日志 / panic 恢复 / 鉴权 拦截器
│   ├── client/       # 客户端拨号助手
│   └── errs/         # gRPC 状态错误封装
├── services/         # 三个服务的业务实现
│   ├── user/
│   ├── product/
│   └── order/
├── cmd/              # 各程序入口
│   ├── user-service/  product-service/  order-service/
│   └── demo-client/
├── docs/guide.md     # 一步步讲解文档
└── Makefile
```

## 快速开始

```bash
# 1. 安装代码生成插件（一次性）
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 2. 拉依赖
make tidy

# 3. （可选）从 proto 重新生成代码
make proto

# 4. 编译
make build

# 5. 后台启动三个服务
make run-all

# 6. 运行演示客户端，观察级联调用与四种 RPC
make demo
```

## 预期输出（节选）

```
================ 1. Unary + 级联调用：CreateOrder ================
   创建成功：id=o1 user=u1 total=31700分 ...
================ 5. Bidirectional streaming ================
   收到更新: order=o1 status=ORDER_STATUS_PENDING  msg="order received"
   收到更新: order=o1 status=ORDER_STATUS_PAID     msg="payment confirmed"
   收到更新: order=o1 status=ORDER_STATUS_SHIPPED  msg="order shipped"
✅ 全部演示完成
```

> 说明：本项目用明文（insecure）传输以便学习；生产环境请改用 TLS。
