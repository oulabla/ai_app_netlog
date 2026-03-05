package service

import (
	"context"

	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

func (s *Service) GetByID(ctx context.Context, id int64) (*datastruct.Netlog, error) {
	return s.repo.GetByID(ctx, id)
}
