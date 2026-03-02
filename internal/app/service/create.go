package service

import (
	"context"

	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

func (s *Service) Create(ctx context.Context, netlog *datastruct.Netlog) (int64, error) {
	return s.repo.Insert(ctx, netlog)
}
