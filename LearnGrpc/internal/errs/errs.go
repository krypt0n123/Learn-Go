// Package errs 封装常用的 gRPC 状态错误。
//
// gRPC 用 status.Error(code, msg) 把错误码与错误信息打包，
// 客户端可以用 status.Code(err) 还原出 codes.* 来做分支处理。
// 统一在这里构造，避免在业务代码里到处 import codes/status。
package errs

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InvalidArgument 参数非法（400）。
func InvalidArgument(format string, a ...any) error {
	return status.Errorf(codes.InvalidArgument, format, a...)
}

// NotFound 资源不存在（404）。
func NotFound(format string, a ...any) error {
	return status.Errorf(codes.NotFound, format, a...)
}

// AlreadyExists 资源已存在（409）。
func AlreadyExists(format string, a ...any) error {
	return status.Errorf(codes.AlreadyExists, format, a...)
}

// FailedPrecondition 前置条件不满足（如库存不足）。
func FailedPrecondition(format string, a ...any) error {
	return status.Errorf(codes.FailedPrecondition, format, a...)
}

// Internal 内部错误（500）。把底层 error 的信息透出，便于学习时排查。
func Internal(err error) error {
	return status.Error(codes.Internal, fmt.Sprintf("internal error: %v", err))
}
