// Code generated. DO NOT EDIT.

package netlog

import (
	"context"

	pb "github.com/oulabla/ai_app_netlog/gen/go/netlog/v1"
)

func (c *Controller) GetNetlog(ctx context.Context, req *pb.GetNetlogRequest) (*pb.GetNetlogResponse, error) {
	nl, err := c.service.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &pb.GetNetlogResponse{
		Netlog: convertNetlog(nl),
	}, nil
}
