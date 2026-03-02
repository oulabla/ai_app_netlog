package service

import (
	"context"

	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

type repository interface {
	Insert(ctx context.Context, netlog *datastruct.Netlog) (int64, error)
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}
