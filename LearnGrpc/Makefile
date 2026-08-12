.PHONY: proto build run-user run-product run-order run-all demo tidy vet clean

PROTOC     := protoc
PROTO_FILES := proto/user/v1/user.proto proto/product/v1/product.proto proto/order/v1/order.proto
BIN        := bin

## 生成 Go 代码（需要 protoc + protoc-gen-go + protoc-gen-go-grpc）
proto:
	$(PROTOC) --proto_path=proto \
	  --go_out=. --go_opt=module=learngrpc \
	  --go-grpc_out=. --go-grpc_opt=module=learngrpc \
	  $(PROTO_FILES)

## 编译所有服务与 demo 客户端到 bin/
build:
	mkdir -p $(BIN)
	go build -o $(BIN)/user-service    ./cmd/user-service
	go build -o $(BIN)/product-service ./cmd/product-service
	go build -o $(BIN)/order-service   ./cmd/order-service
	go build -o $(BIN)/demo-client     ./cmd/demo-client

run-user:
	go run ./cmd/user-service

run-product:
	go run ./cmd/product-service

run-order:
	go run ./cmd/order-service

## 后台启动三个服务（日志写入 /tmp/*.log）
run-all:
	@echo "starting services in background..."
	@go run ./cmd/user-service    > /tmp/user-service.log    2>&1 &
	@go run ./cmd/product-service > /tmp/product-service.log 2>&1 &
	@go run ./cmd/order-service   > /tmp/order-service.log   2>&1 &
	@sleep 1
	@echo "services started. next: make demo   (logs: /tmp/*-service.log)"

## 运行演示客户端（需先 make run-all）
demo:
	go run ./cmd/demo-client

tidy:
	go mod tidy

vet:
	go vet ./...

clean:
	rm -rf $(BIN)
