package service

import (
	"context"

	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

func (s *Service) GetList(ctx context.Context, filter *datastruct.NetlogFilter) ([]*datastruct.Netlog, int64, error) {
	return s.repo.GetList(ctx, filter)
}
