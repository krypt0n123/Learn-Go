// Package product 实现 ProductServiceServer，内存存储商品与库存。
package product

import (
	"context"
	"io"
	"sync"

	productv1 "learngrpc/gen/product/v1"
	"learngrpc/internal/errs"
)

// Service 实现 productv1.ProductServiceServer。
type Service struct {
	productv1.UnimplementedProductServiceServer

	mu       sync.RWMutex
	products map[string]*productv1.Product
}

func New() *Service {
	s := &Service{products: make(map[string]*productv1.Product)}
	s.seed()
	return s
}

func (s *Service) seed() {
	for _, p := range []*productv1.Product{
		{Id: "p1", Name: "Keyboard", PriceCents: 19900, Stock: 10},
		{Id: "p2", Name: "Mouse", PriceCents: 5900, Stock: 25},
		{Id: "p3", Name: "Monitor", PriceCents: 129900, Stock: 5},
	} {
		s.products[p.Id] = p
	}
}

// GetProduct 获取商品（Unary）。
func (s *Service) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.Product, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.products[req.GetId()]
	if !ok {
		return nil, errs.NotFound("product %q not found", req.GetId())
	}
	return p, nil
}

// CheckStock 校验库存是否充足（Unary）。order-service 会调用它。
func (s *Service) CheckStock(ctx context.Context, req *productv1.CheckStockRequest) (*productv1.StockInfo, error) {
	if req.GetQuantity() <= 0 {
		return nil, errs.InvalidArgument("quantity must be positive")
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.products[req.GetProductId()]
	if !ok {
		return nil, errs.NotFound("product %q not found", req.GetProductId())
	}
	available := p.GetStock() >= req.GetQuantity()
	return &productv1.StockInfo{
		ProductId: p.GetId(),
		Available: available,
		Remaining: p.GetStock(),
	}, nil
}

// ListProducts 列出商品（Server streaming）。
func (s *Service) ListProducts(req *productv1.ListProductsRequest, stream productv1.ProductService_ListProductsServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	pageSize := int(req.GetPageSize())
	sent := 0
	for _, p := range s.products {
		if pageSize > 0 && sent >= pageSize {
			break
		}
		if err := stream.Send(p); err != nil {
			return err
		}
		sent++
	}
	return nil
}

// AdjustStock 批量调整库存（Client streaming）。
// 客户端连续发送多条 StockAdjustment，服务端边收边改，
// 最后返回一条汇总。Recv 返回 io.EOF 表示客户端发完。
func (s *Service) AdjustStock(stream productv1.ProductService_AdjustStockServer) error {
	processed := 0
	for {
		adj, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&productv1.AdjustStockSummary{
				Processed: int32(processed),
				Ok:        true,
			})
		}
		if err != nil {
			return err
		}
		if adj.GetDelta() == 0 {
			continue
		}
		s.mu.Lock()
		p, ok := s.products[adj.GetProductId()]
		if !ok {
			s.mu.Unlock()
			return errs.NotFound("product %q not found", adj.GetProductId())
		}
		newStock := p.GetStock() + adj.GetDelta()
		if newStock < 0 {
			s.mu.Unlock()
			return errs.FailedPrecondition("insufficient stock for %q", adj.GetProductId())
		}
		p.Stock = newStock
		s.mu.Unlock()
		processed++
	}
}
