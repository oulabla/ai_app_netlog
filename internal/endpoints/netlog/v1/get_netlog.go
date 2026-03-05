// Code generated. DO NOT EDIT.

package netlog

import (
	"context"

	pb "github.com/oulabla/ai_app_netlog/gen/go/netlog/v1"
	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

func (c *Controller) GetNetlog(ctx context.Context, req *pb.GetNetlogRequest) (*pb.GetNetlogResponse, error) {
	nl, err := c.service.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	if nl == nil {
		return nil, datastruct.NewNotFound("netlog not found")
	}

	return &pb.GetNetlogResponse{
		Netlog: convertNetlog(nl),
	}, nil
}
