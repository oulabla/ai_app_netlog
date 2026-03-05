// Code generated. DO NOT EDIT.

package netlog

import (
	"context"
	"time"

	pb "github.com/oulabla/ai_app_netlog/gen/go/netlog/v1"
	"github.com/oulabla/ai_app_netlog/internal/datastruct"
)

func (c *Controller) CreateNetlog(ctx context.Context, req *pb.CreateNetlogRequest) (*pb.CreateNetlogResponse, error) {
	netlogErr := req.GetNetlog().GetError()
	netlog := &datastruct.Netlog{
		CreatedAt:            time.Now(),
		ClientID:             req.GetNetlog().GetClientId(),
		AppName:              req.Netlog.GetAppName(),
		Keywords:             req.GetNetlog().GetKeywords(),
		Parameters:           req.GetNetlog().GetParameters().AsMap(),
		Error:                &netlogErr,
		NumBeforeAiFilter:    int(req.GetNetlog().GetNumBeforeAiFilter()),
		NumAfterAiFilter:     int(req.GetNetlog().GetNumAfterAiFilter()),
		ResultBeforeAiFilter: req.GetNetlog().GetResultBeforeAiFilter().AsMap(),
		Result:               req.GetNetlog().GetResult().AsMap(),
	}

	id, err := c.service.Create(ctx, netlog)
	if err != nil {
		return nil, err
	}

	return &pb.CreateNetlogResponse{
		Id: id,
	}, nil
}
