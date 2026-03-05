// internal/server/error_interceptor.go
package server

import (
	"context"
	"errors"

	"github.com/oulabla/ai_app_netlog/internal/datastruct"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UnaryErrorInterceptor перехватывает ошибки и преобразует AppError в соответствующий gRPC status
func UnaryErrorInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err == nil {
		return resp, nil
	}

	// Пытаемся распаковать как *AppError
	var appErr *datastruct.AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case datastruct.CodeNotFound:
			return nil, status.Error(codes.NotFound, appErr.Error())

		case datastruct.CodeValidationError:
			return nil, status.Error(codes.InvalidArgument, appErr.Error())

		case datastruct.CodeUnauthorized:
			return nil, status.Error(codes.Unauthenticated, appErr.Error())

		case datastruct.CodeConflict:
			return nil, status.Error(codes.AlreadyExists, appErr.Error())

		case datastruct.CodeInternalError:
			return nil, status.Error(codes.Internal, appErr.Error())

		case datastruct.CodePermissionDenied:
			return nil, status.Error(codes.PermissionDenied, appErr.Error())

		case datastruct.CodeRateLimited:
			return nil, status.Error(codes.ResourceExhausted, appErr.Error())

		default:
			// Неизвестный код → всё равно Internal, но можно логировать
			return nil, status.Error(codes.Internal, appErr.Error())
		}
	}

	// Если это не AppError — оставляем как есть (или тоже Internal)
	// Можно здесь добавить логирование неизвестных ошибок, если нужно
	return nil, status.Error(codes.Unknown, err.Error())
}
