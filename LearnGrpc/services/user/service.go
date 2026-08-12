// Package user 实现 UserServiceServer，用一个内存 map 模拟数据库。
package user

import (
	"context"
	"fmt"
	"sync"
	"time"

	userv1 "learngrpc/gen/user/v1"
	"learngrpc/internal/errs"
)

// Service 实现 userv1.UserServiceServer。
type Service struct {
	userv1.UnimplementedUserServiceServer

	mu    sync.RWMutex
	users map[string]*userv1.User
}

// New 创建用户服务并预置一些示例数据。
func New() *Service {
	s := &Service{users: make(map[string]*userv1.User)}
	s.seed()
	return s
}

func (s *Service) seed() {
	now := time.Now().Unix()
	for _, u := range []*userv1.User{
		{Id: "u1", Name: "Alice", Email: "alice@example.com", CreatedAt: now},
		{Id: "u2", Name: "Bob", Email: "bob@example.com", CreatedAt: now},
		{Id: "u3", Name: "Carol", Email: "carol@example.com", CreatedAt: now},
	} {
		s.users[u.Id] = u
	}
}

// CreateUser 创建用户（Unary）。
func (s *Service) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.User, error) {
	if req.GetName() == "" || req.GetEmail() == "" {
		return nil, errs.InvalidArgument("name and email are required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, u := range s.users {
		if u.Email == req.GetEmail() {
			return nil, errs.AlreadyExists("user with email %q already exists", req.GetEmail())
		}
	}
	u := &userv1.User{
		Id:        fmt.Sprintf("u%d", len(s.users)+1),
		Name:      req.GetName(),
		Email:     req.GetEmail(),
		CreatedAt: time.Now().Unix(),
	}
	s.users[u.Id] = u
	return u, nil
}

// GetUser 根据 id 获取用户（Unary）。
func (s *Service) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[req.GetId()]
	if !ok {
		return nil, errs.NotFound("user %q not found", req.GetId())
	}
	return u, nil
}

// ListUsers 列出用户（Server streaming）。
// 服务端通过 stream.Send 逐条把 User 推给客户端，直到全部发完。
func (s *Service) ListUsers(req *userv1.ListUsersRequest, stream userv1.UserService_ListUsersServer) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	pageSize := int(req.GetPageSize())
	sent := 0
	for _, u := range s.users {
		if pageSize > 0 && sent >= pageSize {
			break
		}
		if err := stream.Send(u); err != nil {
			return err // 通常是客户端断开
		}
		sent++
	}
	return nil
}