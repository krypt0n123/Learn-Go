// Package order 实现 OrderServiceServer。
//
// 它是整个示例的核心：order-service 既是 gRPC 服务端（对外提供订单 API），
// 又是 gRPC 客户端——在处理 CreateOrder 时会去调用 user-service.GetUser
// 和 product-service.CheckStock，演示“服务间调用”。
package order

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	orderv1 "learngrpc/gen/order/v1"
	productv1 "learngrpc/gen/product/v1"
	userv1 "learngrpc/gen/user/v1"
	"learngrpc/internal/errs"
)

// Clients 聚合 order-service 需要调用的下游服务客户端。
type Clients struct {
	User    userv1.UserServiceClient
	Product productv1.ProductServiceClient
}

// Service 实现 orderv1.OrderServiceServer。
type Service struct {
	orderv1.UnimplementedOrderServiceServer

	clients Clients

	mu     sync.RWMutex
	orders map[string]*orderv1.Order
	nextID int
}

func New(clients Clients) *Service {
	return &Service{
		clients: clients,
		orders:  make(map[string]*orderv1.Order),
	}
}

// CreateOrder 创建订单（Unary）。
//
// 调用链（核心演示）：
//  1. 调用 user-service.GetUser 校验用户存在
//  2. 逐个调用 product-service.CheckStock 校验库存
//  3. 调用 product-service.GetProduct 取真实价格，避免被客户端伪造
//  4. 全部通过后才生成订单
//
// 注意：下游调用复用入参 ctx，这样上游设置的超时会沿调用链传播。
func (s *Service) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.Order, error) {
	if req.GetUserId() == "" {
		return nil, errs.InvalidArgument("user_id is required")
	}
	if len(req.GetItems()) == 0 {
		return nil, errs.InvalidArgument("items must not be empty")
	}

	// 1) 校验用户：跨服务调用 user-service
	if _, err := s.clients.User.GetUser(ctx, &userv1.GetUserRequest{Id: req.GetUserId()}); err != nil {
		// 把下游的 gRPC 状态码透传/降级为更合适的语义
		if status.Code(err) == codes.NotFound {
			return nil, errs.FailedPrecondition("user %q does not exist", req.GetUserId())
		}
		return nil, errs.Internal(fmt.Errorf("user-service: %w", err))
	}

	// 2) + 3) 逐个校验库存并取真实价格
	items := make([]*orderv1.OrderItem, 0, len(req.GetItems()))
	totalCents := int64(0)
	for _, it := range req.GetItems() {
		stock, err := s.clients.Product.CheckStock(ctx, &productv1.CheckStockRequest{
			ProductId: it.GetProductId(),
			Quantity:  it.GetQuantity(),
		})
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return nil, errs.FailedPrecondition("product %q does not exist", it.GetProductId())
			}
			return nil, errs.Internal(fmt.Errorf("product-service CheckStock: %w", err))
		}
		if !stock.GetAvailable() {
			return nil, errs.FailedPrecondition("insufficient stock for product %q (remaining %d)",
				it.GetProductId(), stock.GetRemaining())
		}

		prod, err := s.clients.Product.GetProduct(ctx, &productv1.GetProductRequest{Id: it.GetProductId()})
		if err != nil {
			return nil, errs.Internal(fmt.Errorf("product-service GetProduct: %w", err))
		}
		items = append(items, &orderv1.OrderItem{
			ProductId:  prod.GetId(),
			Quantity:   it.GetQuantity(),
			PriceCents: prod.GetPriceCents(),
		})
		totalCents += prod.GetPriceCents() * int64(it.GetQuantity())
	}

	// 4) 生成订单
	s.mu.Lock()
	s.nextID++
	id := fmt.Sprintf("o%d", s.nextID)
	o := &orderv1.Order{
		Id:         id,
		UserId:     req.GetUserId(),
		Items:      items,
		TotalCents: totalCents,
		Status:     orderv1.OrderStatus_ORDER_STATUS_PENDING,
		CreatedAt:  time.Now().Unix(),
	}
	s.orders[id] = o
	s.mu.Unlock()

	return o, nil
}

// GetOrder 查询订单（Unary）。
func (s *Service) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.Order, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.orders[req.GetId()]
	if !ok {
		return nil, errs.NotFound("order %q not found", req.GetId())
	}
	return o, nil
}

// StreamOrderUpdates 订单状态更新（Bidirectional streaming）。
// 客户端发一条 OrderQuery，服务端就把该订单的“模拟状态流转”逐条回推，
// 双方都可以主动结束流。
func (s *Service) StreamOrderUpdates(stream orderv1.OrderService_StreamOrderUpdatesServer) error {
	transitions := []struct {
		status  orderv1.OrderStatus
		message string
	}{
		{orderv1.OrderStatus_ORDER_STATUS_PENDING, "order received"},
		{orderv1.OrderStatus_ORDER_STATUS_PAID, "payment confirmed"},
		{orderv1.OrderStatus_ORDER_STATUS_SHIPPED, "order shipped"},
	}
	for {
		q, err := stream.Recv()
		if err == io.EOF {
			return nil // 客户端结束，我们也结束
		}
		if err != nil {
			return err
		}
		// 校验订单是否存在
		s.mu.RLock()
		_, exists := s.orders[q.GetOrderId()]
		s.mu.RUnlock()
		if !exists {
			if err := stream.Send(&orderv1.OrderUpdate{
				OrderId: q.GetOrderId(),
				Message: "order not found",
			}); err != nil {
				return err
			}
			continue
		}
		// 推送模拟的状态流转
		for _, t := range transitions {
			if err := stream.Send(&orderv1.OrderUpdate{
				OrderId: q.GetOrderId(),
				Status:  t.status,
				Message: t.message,
			}); err != nil {
				return err
			}
			time.Sleep(200 * time.Millisecond) // 仅用于演示节奏
		}
	}
}
