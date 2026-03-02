// Code generated. DO NOT EDIT.

package netlog

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/oulabla/ai_app_netlog/gen/go/netlog/v1"
	main_service "github.com/oulabla/ai_app_netlog/internal/app/service"
	"github.com/oulabla/ai_app_netlog/internal/config"
	"github.com/oulabla/ai_app_netlog/internal/datastruct"
	"github.com/oulabla/ai_app_netlog/internal/server"
	"github.com/rs/zerolog/log"
)

type service interface {
	Create(ctx context.Context, netlog *datastruct.Netlog) (int64, error)
}

type Controller struct {
	pb.UnimplementedNetlogServiceServer
	// сюда внедряем зависимости (usecases, logger, repositories и т.д.)
	// usecase usecase.UserServiceUsecase
	// logger  *zap.Logger
	service service
}

func NewController() *Controller {
	s, err := server.GetWithType[*main_service.Service](config.AppService)
	if err != nil {
		log.Fatal().Msg("di service")
	}
	return &Controller{
		service: s,
	}
}

// Register регистрирует сервис в grpc-сервере
func (c *Controller) Register(srv *grpc.Server) {
	pb.RegisterNetlogServiceServer(srv, c)
	reflection.Register(srv)
}
