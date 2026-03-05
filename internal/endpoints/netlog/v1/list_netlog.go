// Code generated. DO NOT EDIT.

package netlog

import (
	"context"

	pb "github.com/oulabla/ai_app_netlog/gen/go/netlog/v1"
	"github.com/oulabla/ai_app_netlog/internal/datastruct"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Controller) ListNetlog(ctx context.Context, req *pb.ListNetlogRequest) (*pb.ListNetlogResponse, error) {
	filter := сonvertListNetlogRequestToFilter(req)

	logs, lastID, err := c.service.GetList(ctx, filter)
	if err != nil {
		return nil, err
	}

	// TODO: implement logic
	return &pb.ListNetlogResponse{
		Items:      convertNetlogs(logs),
		NextLastId: &lastID,
	}, nil
}

func сonvertListNetlogRequestToFilter(req *pb.ListNetlogRequest) *datastruct.NetlogFilter {
	if req == nil {
		return &datastruct.NetlogFilter{
			Limit: 20, // значение по умолчанию
		}
	}

	filter := &datastruct.NetlogFilter{
		Limit:  int(req.Limit),
		LastID: req.LastId,
	}

	// Копирование строковых полей
	if req.ClientId != nil {
		filter.ClientID = *req.ClientId
	}

	if req.AppName != nil {
		filter.AppName = req.AppName
	}

	// Копирование слайса Keywords
	if len(req.Keywords) > 0 {
		filter.Keywords = make([]string, len(req.Keywords))
		copy(filter.Keywords, req.Keywords)
	}

	// Копирование булевых полей
	if req.HasError != nil {
		filter.HasError = req.HasError
	}

	// Копирование числовых полей (int32 -> int)
	if req.MinBeforeAi != nil {
		val := int(*req.MinBeforeAi)
		filter.MinBeforeAI = &val
	}

	if req.MaxBeforeAi != nil {
		val := int(*req.MaxBeforeAi)
		filter.MaxBeforeAI = &val
	}

	// Конвертация Timestamp в time.Time
	if req.FromTime != nil {
		fromTime := req.FromTime.AsTime()
		filter.FromTime = &fromTime
	}

	if req.ToTime != nil {
		toTime := req.ToTime.AsTime()
		filter.ToTime = &toTime
	}

	return filter
}

func convertNetlogs(logs []*datastruct.Netlog) []*pb.Netlog {
	res := make([]*pb.Netlog, 0, len(logs))
	for _, l := range logs {
		res = append(res, convertNetlog(l))
	}
	return res
}

func convertNetlog(l *datastruct.Netlog) *pb.Netlog {
	params, err := structpb.NewStruct(l.Parameters)
	if err != nil {
		log.Warn().Err(err).Msg("list_netlog: convert parameters")
	}

	result, err := structpb.NewStruct(l.Result)
	if err != nil {
		log.Warn().Err(err).Msg("list_netlog: convert result")
	}

	resultBeforeAI, err := structpb.NewStruct(l.ResultBeforeAiFilter)
	if err != nil {
		log.Warn().Err(err).Msg("list_netlog: convert result")
	}

	return &pb.Netlog{
		Id:                   l.ID,
		Keywords:             l.Keywords,
		Parameters:           params,
		Error:                l.Error,
		NumBeforeAiFilter:    int32(l.NumBeforeAiFilter),
		NumAfterAiFilter:     int32(l.NumAfterAiFilter),
		ClientId:             l.ClientID,
		AppName:              l.AppName,
		CreatedAt:            timestamppb.New(l.CreatedAt),
		Result:               result,
		ResultBeforeAiFilter: resultBeforeAI,
	}
}
