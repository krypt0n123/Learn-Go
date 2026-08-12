// demo-client 演示如何作为 gRPC 客户端调用各个服务。
//
// 它会依次演示：
//  1. Unary + 服务间级联：调用 order-service.CreateOrder，
//     该调用内部又会触发 order-service -> user-service / product-service 的级联
//  2. Unary：order-service.GetOrder
//  3. Server streaming：直连 product-service.ListProducts
//  4. Client streaming：直连 product-service.AdjustStock
//  5. Bidirectional streaming：order-service.StreamOrderUpdates
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/status"

	orderv1 "learngrpc/gen/order/v1"
	productv1 "learngrpc/gen/product/v1"
	"learngrpc/internal/client"
	"learngrpc/internal/server"
)

func main() {
	log.SetFlags(log.Ltime)

	// order-service 对外需要鉴权 token；product-service 是内部服务，不需要。
	orderConn, err := client.Dial(server.EnvOr("ORDER_SERVICE_ADDR", "localhost:50053"), "dev-secret")
	if err != nil {
		log.Fatalf("dial order-service: %v", err)
	}
	defer orderConn.Close()
	productConn, err := client.Dial(server.EnvOr("PRODUCT_SERVICE_ADDR", "localhost:50052"), "")
	if err != nil {
		log.Fatalf("dial product-service: %v", err)
	}
	defer productConn.Close()

	orderClient := orderv1.NewOrderServiceClient(orderConn)
	productClient := productv1.NewProductServiceClient(productConn)

	banner("1. Unary + 级联调用：CreateOrder（会触发 order -> user / product）")
	// 给整个级联调用设置 3s 超时；该 deadline 会沿调用链传播到下游服务。
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	orderID := demoCreateOrder(ctx, orderClient)
	cancel()

	banner("2. Unary：GetOrder")
	demoGetOrder(context.Background(), orderClient, orderID)

	banner("3. Server streaming：product-service.ListProducts")
	demoListProducts(context.Background(), productClient)

	banner("4. Client streaming：product-service.AdjustStock")
	demoAdjustStock(context.Background(), productClient)

	banner("5. Bidirectional streaming：order-service.StreamOrderUpdates")
	demoStreamOrderUpdates(context.Background(), orderClient, orderID)

	fmt.Println("\n✅ 全部演示完成")
}

// ---------- 1+2. Unary ----------

func demoCreateOrder(ctx context.Context, c orderv1.OrderServiceClient) string {
	req := &orderv1.CreateOrderRequest{
		UserId: "u1",
		Items: []*orderv1.OrderItem{
			{ProductId: "p1", Quantity: 1},
			{ProductId: "p2", Quantity: 2},
		},
	}
	o, err := c.CreateOrder(ctx, req)
	if err != nil {
		log.Fatalf("CreateOrder failed: code=%v %v", status.Code(err), err)
	}
	fmt.Printf("   创建成功：%s\n", formatOrder(o))
	return o.GetId()
}

func demoGetOrder(ctx context.Context, c orderv1.OrderServiceClient, id string) {
	o, err := c.GetOrder(ctx, &orderv1.GetOrderRequest{Id: id})
	if err != nil {
		log.Fatalf("GetOrder failed: code=%v %v", status.Code(err), err)
	}
	fmt.Printf("   查询结果：%s\n", formatOrder(o))
}

// ---------- 3. Server streaming ----------

func demoListProducts(ctx context.Context, c productv1.ProductServiceClient) {
	stream, err := c.ListProducts(ctx, &productv1.ListProductsRequest{PageSize: 10})
	if err != nil {
		log.Fatalf("ListProducts open stream: %v", err)
	}
	for {
		p, err := stream.Recv()
		if err == io.EOF {
			break // 流结束
		}
		if err != nil {
			log.Fatalf("ListProducts recv: %v", err)
		}
		fmt.Printf("   收到商品: id=%s name=%s price=%d.%02d元 stock=%d\n",
			p.GetId(), p.GetName(), p.GetPriceCents()/100, p.GetPriceCents()%100, p.GetStock())
	}
}

// ---------- 4. Client streaming ----------

func demoAdjustStock(ctx context.Context, c productv1.ProductServiceClient) {
	stream, err := c.AdjustStock(ctx)
	if err != nil {
		log.Fatalf("AdjustStock open stream: %v", err)
	}
	// 连续发送多条库存调整
	adjustments := []*productv1.StockAdjustment{
		{ProductId: "p1", Delta: 5},  // 入库 5
		{ProductId: "p2", Delta: -3}, // 出库 3
		{ProductId: "p3", Delta: 2},  // 入库 2
	}
	for _, a := range adjustments {
		if err := stream.Send(a); err != nil {
			log.Fatalf("AdjustStock send: %v", err)
		}
		fmt.Printf("   已发送调整: product=%s delta=%+d\n", a.GetProductId(), a.GetDelta())
	}
	summary, err := stream.CloseAndRecv() // 发完，等汇总
	if err != nil {
		log.Fatalf("AdjustStock close: %v", err)
	}
	fmt.Printf("   服务端汇总: processed=%d ok=%v\n", summary.GetProcessed(), summary.GetOk())
}

// ---------- 5. Bidirectional streaming ----------

func demoStreamOrderUpdates(ctx context.Context, c orderv1.OrderServiceClient, orderID string) {
	stream, err := c.StreamOrderUpdates(ctx)
	if err != nil {
		log.Fatalf("StreamOrderUpdates open stream: %v", err)
	}
	// 先发一条查询
	if err := stream.Send(&orderv1.OrderQuery{OrderId: orderID}); err != nil {
		log.Fatalf("StreamOrderUpdates send: %v", err)
	}
	fmt.Printf("   已发送查询: order=%s，等待状态推送...\n", orderID)
	// 接收服务端推送的状态更新；收到“已发货”后主动 CloseSend，
	// 服务端的 Recv 会收到 io.EOF 从而结束流，客户端再排空即可。
	closed := false
	for {
		upd, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("StreamOrderUpdates recv: %v", err)
		}
		fmt.Printf("   收到更新: order=%s status=%s msg=%q\n",
			upd.GetOrderId(), upd.GetStatus(), upd.GetMessage())
		if upd.GetStatus() == orderv1.OrderStatus_ORDER_STATUS_SHIPPED && !closed {
			closed = true
			if err := stream.CloseSend(); err != nil {
				log.Fatalf("StreamOrderUpdates close send: %v", err)
			}
		}
	}
}

// ---------- 辅助 ----------

func banner(title string) {
	fmt.Printf("\n================ %s ================\n", title)
}

func formatOrder(o *orderv1.Order) string {
	items := ""
	for i, it := range o.GetItems() {
		if i > 0 {
			items += ", "
		}
		items += fmt.Sprintf("%s x%d@%d分", it.GetProductId(), it.GetQuantity(), it.GetPriceCents())
	}
	return fmt.Sprintf("id=%s user=%s total=%d分 status=%s items=[%s]",
		o.GetId(), o.GetUserId(), o.GetTotalCents(), o.GetStatus(), items)
}
