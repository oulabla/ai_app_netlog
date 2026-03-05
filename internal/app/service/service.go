package service

import (
	"context"

	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

type repository interface {
	Insert(ctx context.Context, netlog *datastruct.Netlog) (int64, error)
	GetList(ctx context.Context, filter *datastruct.NetlogFilter) ([]*datastruct.Netlog, int64, error)
	GetByID(ctx context.Context, id int64) (*datastruct.Netlog, error)
}

type Service struct {
	repo repository
}

func NewService(repo repository) *Service {
	return &Service{
		repo: repo,
	}
}
