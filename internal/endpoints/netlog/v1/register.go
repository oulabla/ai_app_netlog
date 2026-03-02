// Code generated. DO NOT EDIT.

package netlog

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	pb "github.com/oulabla/ai_app_netlog/gen/go/netlog/v1"
	"github.com/oulabla/ai_app_netlog/internal/server"
)

func init() {
	server.RegisterGRPC(func(srv *grpc.Server) {
		pb.RegisterNetlogServiceServer(srv, NewController())
	})

	server.RegisterGateway(func(ctx context.Context, mux *runtime.ServeMux, grpcAddr string, opts []grpc.DialOption) error {
		return pb.RegisterNetlogServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	})

	server.RegisterSwagger(func() server.SwaggerConfig {
		return server.SwaggerConfig{
			FileName: "all-apis.swagger.json", // фиксированное имя
		}
	})
}
