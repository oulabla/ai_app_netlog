package service

import (
	"context"

	"github.com/oulabla/ai_app_netlog/internal/datastruct"
	"github.com/oulabla/ai_app_netlog/internal/metric"
)

func (s *Service) Create(ctx context.Context, netlog *datastruct.Netlog) (int64, error) {
	id, err := s.repo.Insert(ctx, netlog)
	if err != nil {
		return 0, err
	}

	metric.IncNetlogCreated(netlog.AppName)

	return id, nil
}
