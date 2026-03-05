package netlog

import (
	"context"
	"errors"
	"testing"
	"time"

	pb "github.com/oulabla/ai_app_netlog/gen/go/netlog/v1"
	"github.com/oulabla/ai_app_netlog/internal/datastruct"
	"github.com/oulabla/ai_app_netlog/internal/endpoints/netlog/v1/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestController_ListNetlog(t *testing.T) {
	now := time.Now().UTC()
	errMsg := "error message"
	nextLastID := int64(50)

	tests := []struct {
		name          string
		req           *pb.ListNetlogRequest
		mockSetup     func(mockService *mocks.MockService)
		expectedResp  *pb.ListNetlogResponse
		expectedError error
	}{
		{
			name: "successful list with all filters",
			req: &pb.ListNetlogRequest{
				Limit:       20,
				LastId:      ptrInt64(100),
				ClientId:    ptrString("client-123"),
				AppName:     ptrString("test-app"),
				Keywords:    []string{"kw1", "kw2"},
				HasError:    ptrBool(true),
				MinBeforeAi: ptrInt32(10),
				MaxBeforeAi: ptrInt32(100),
				FromTime:    timestamppb.New(now.Add(-24 * time.Hour)),
				ToTime:      timestamppb.New(now),
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetList", mock.Anything, mock.MatchedBy(func(filter *datastruct.NetlogFilter) bool {
					return filter.Limit == 20 &&
						*filter.LastID == 100 &&
						filter.ClientID == "client-123" &&
						*filter.AppName == "test-app" &&
						len(filter.Keywords) == 2 &&
						filter.Keywords[0] == "kw1" &&
						filter.Keywords[1] == "kw2" &&
						*filter.HasError == true &&
						*filter.MinBeforeAI == 10 &&
						*filter.MaxBeforeAI == 100 &&
						filter.FromTime != nil &&
						filter.ToTime != nil
				})).Return([]*datastruct.Netlog{
					{
						ID:        123,
						CreatedAt: now,
						ClientID:  "client-123",
						AppName:   "test-app",
						Keywords:  []string{"kw1", "kw2"},
						Parameters: map[string]interface{}{
							"param1": "value1",
						},
						Error:             &errMsg,
						NumBeforeAiFilter: 50,
						NumAfterAiFilter:  25,
						ResultBeforeAiFilter: map[string]interface{}{
							"before": "result",
						},
						Result: map[string]interface{}{
							"after": true,
						},
					},
				}, nextLastID, nil)
			},
			expectedResp: &pb.ListNetlogResponse{
				Items: []*pb.Netlog{
					{
						Id:        123,
						CreatedAt: timestamppb.New(now),
						ClientId:  "client-123",
						AppName:   "test-app",
						Keywords:  []string{"kw1", "kw2"},
						Parameters: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"param1": structpb.NewStringValue("value1"),
							},
						},
						Error:             &errMsg,
						NumBeforeAiFilter: 50,
						NumAfterAiFilter:  25,
						ResultBeforeAiFilter: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"before": structpb.NewStringValue("result"),
							},
						},
						Result: &structpb.Struct{
							Fields: map[string]*structpb.Value{
								"after": structpb.NewBoolValue(true),
							},
						},
					},
				},
				NextLastId: &nextLastID,
			},
			expectedError: nil,
		},
		{
			name: "successful list with default filter (nil request)",
			req:  nil,
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetList", mock.Anything, mock.MatchedBy(func(filter *datastruct.NetlogFilter) bool {
					return filter.Limit == 20 && // значение по умолчанию
						filter.LastID == nil &&
						filter.ClientID == "" &&
						filter.AppName == nil &&
						len(filter.Keywords) == 0 &&
						filter.HasError == nil &&
						filter.MinBeforeAI == nil &&
						filter.MaxBeforeAI == nil &&
						filter.FromTime == nil &&
						filter.ToTime == nil
				})).Return([]*datastruct.Netlog{}, int64(0), nil)
			},
			expectedResp: &pb.ListNetlogResponse{
				Items:      []*pb.Netlog{},
				NextLastId: ptrInt64(0),
			},
			expectedError: nil,
		},
		{
			name: "successful list with minimal fields",
			req: &pb.ListNetlogRequest{
				Limit: 30,
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetList", mock.Anything, mock.MatchedBy(func(filter *datastruct.NetlogFilter) bool {
					return filter.Limit == 30 &&
						filter.LastID == nil &&
						filter.ClientID == "" &&
						filter.AppName == nil &&
						len(filter.Keywords) == 0
				})).Return([]*datastruct.Netlog{
					{
						ID:                   456,
						CreatedAt:            now,
						ClientID:             "client-456",
						AppName:              "minimal-app",
						Keywords:             []string{},
						Parameters:           map[string]interface{}{},
						Error:                nil,
						NumBeforeAiFilter:    0,
						NumAfterAiFilter:     0,
						ResultBeforeAiFilter: map[string]interface{}{},
						Result:               map[string]interface{}{},
					},
				}, int64(45), nil)
			},
			expectedResp: &pb.ListNetlogResponse{
				Items: []*pb.Netlog{
					{
						Id:                   456,
						CreatedAt:            timestamppb.New(now),
						ClientId:             "client-456",
						AppName:              "minimal-app",
						Keywords:             []string{},
						Parameters:           &structpb.Struct{},
						Error:                nil,
						NumBeforeAiFilter:    0,
						NumAfterAiFilter:     0,
						ResultBeforeAiFilter: &structpb.Struct{},
						Result:               &structpb.Struct{},
					},
				},
				NextLastId: ptrInt64(45),
			},
			expectedError: nil,
		},
		{
			name: "service error",
			req: &pb.ListNetlogRequest{
				Limit: 20,
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetList", mock.Anything, mock.Anything).
					Return(nil, int64(0), errors.New("database error"))
			},
			expectedResp:  nil,
			expectedError: errors.New("database error"),
		},
		{
			name: "empty result",
			req: &pb.ListNetlogRequest{
				Limit:    50,
				ClientId: ptrString("non-existent"),
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetList", mock.Anything, mock.Anything).
					Return([]*datastruct.Netlog{}, int64(0), nil)
			},
			expectedResp: &pb.ListNetlogResponse{
				Items:      []*pb.Netlog{},
				NextLastId: ptrInt64(0),
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(mocks.MockService)

			// Setup mock expectations
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}

			// Create controller with mock service
			controller := &Controller{
				service: mockService,
			}

			// Execute the method
			resp, err := controller.ListNetlog(context.Background(), tt.req)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, len(tt.expectedResp.Items), len(resp.Items))
				assert.Equal(t, tt.expectedResp.NextLastId, resp.NextLastId)

				// Проверяем каждый элемент, если есть
				for i, expectedItem := range tt.expectedResp.Items {
					actualItem := resp.Items[i]
					assert.Equal(t, expectedItem.Id, actualItem.Id)
					assert.Equal(t, expectedItem.ClientId, actualItem.ClientId)
					assert.Equal(t, expectedItem.AppName, actualItem.AppName)
					assert.Equal(t, expectedItem.Keywords, actualItem.Keywords)
					assert.Equal(t, expectedItem.Error, actualItem.Error)
					assert.Equal(t, expectedItem.NumBeforeAiFilter, actualItem.NumBeforeAiFilter)
					assert.Equal(t, expectedItem.NumAfterAiFilter, actualItem.NumAfterAiFilter)

					// Проверяем Parameters
					if expectedItem.Parameters != nil {
						assert.Equal(t, expectedItem.Parameters.AsMap(), actualItem.Parameters.AsMap())
					}

					// Проверяем ResultBeforeAiFilter
					if expectedItem.ResultBeforeAiFilter != nil {
						assert.Equal(t, expectedItem.ResultBeforeAiFilter.AsMap(), actualItem.ResultBeforeAiFilter.AsMap())
					}

					// Проверяем Result
					if expectedItem.Result != nil {
						assert.Equal(t, expectedItem.Result.AsMap(), actualItem.Result.AsMap())
					}
				}
			}

			// Verify that all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

// Тест для функции конвертации фильтра
func TestConvertListNetlogRequestToFilter(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		req      *pb.ListNetlogRequest
		expected *datastruct.NetlogFilter
	}{
		{
			name: "nil request",
			req:  nil,
			expected: &datastruct.NetlogFilter{
				Limit: 20, // default value
			},
		},
		{
			name: "all fields populated",
			req: &pb.ListNetlogRequest{
				Limit:       50,
				LastId:      ptrInt64(100),
				ClientId:    ptrString("client-123"),
				AppName:     ptrString("app-name"),
				Keywords:    []string{"kw1", "kw2"},
				HasError:    ptrBool(true),
				MinBeforeAi: ptrInt32(10),
				MaxBeforeAi: ptrInt32(100),
				FromTime:    timestamppb.New(now.Add(-24 * time.Hour)),
				ToTime:      timestamppb.New(now),
			},
			expected: &datastruct.NetlogFilter{
				Limit:       50,
				LastID:      ptrInt64(100),
				ClientID:    "client-123",
				AppName:     ptrString("app-name"),
				Keywords:    []string{"kw1", "kw2"},
				HasError:    ptrBool(true),
				MinBeforeAI: ptrInt(10),
				MaxBeforeAI: ptrInt(100),
				FromTime:    ptrTime(now.Add(-24 * time.Hour)),
				ToTime:      ptrTime(now),
			},
		},
		{
			name: "only limit",
			req: &pb.ListNetlogRequest{
				Limit: 30,
			},
			expected: &datastruct.NetlogFilter{
				Limit:  30,
				LastID: nil,
			},
		},
		{
			name: "only last_id",
			req: &pb.ListNetlogRequest{
				LastId: ptrInt64(200),
			},
			expected: &datastruct.NetlogFilter{
				Limit:  0, // zero value
				LastID: ptrInt64(200),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := сonvertListNetlogRequestToFilter(tt.req)

			assert.Equal(t, tt.expected.Limit, result.Limit)
			assert.Equal(t, tt.expected.LastID, result.LastID)
			assert.Equal(t, tt.expected.ClientID, result.ClientID)
			assert.Equal(t, tt.expected.AppName, result.AppName)
			assert.Equal(t, tt.expected.Keywords, result.Keywords)
			assert.Equal(t, tt.expected.HasError, result.HasError)
			assert.Equal(t, tt.expected.MinBeforeAI, result.MinBeforeAI)
			assert.Equal(t, tt.expected.MaxBeforeAI, result.MaxBeforeAI)

			// Проверяем временные поля
			if tt.expected.FromTime != nil {
				assert.Equal(t, tt.expected.FromTime.Unix(), result.FromTime.Unix())
			} else {
				assert.Nil(t, result.FromTime)
			}

			if tt.expected.ToTime != nil {
				assert.Equal(t, tt.expected.ToTime.Unix(), result.ToTime.Unix())
			} else {
				assert.Nil(t, result.ToTime)
			}
		})
	}
}

// Тест для функции конвертации списка netlog
func TestConvertNetlogs(t *testing.T) {
	now := time.Now().UTC()
	errMsg := "error"

	logs := []*datastruct.Netlog{
		{
			ID:        1,
			CreatedAt: now,
			ClientID:  "client-1",
			AppName:   "app-1",
			Keywords:  []string{"kw1"},
			Parameters: map[string]interface{}{
				"param": "value",
			},
			Error:             &errMsg,
			NumBeforeAiFilter: 10,
			NumAfterAiFilter:  5,
			ResultBeforeAiFilter: map[string]interface{}{
				"before": "result",
			},
			Result: map[string]interface{}{
				"after": true,
			},
		},
		{
			ID:                   2,
			CreatedAt:            now,
			ClientID:             "client-2",
			AppName:              "app-2",
			Keywords:             []string{},
			Parameters:           map[string]interface{}{},
			Error:                nil,
			NumBeforeAiFilter:    0,
			NumAfterAiFilter:     0,
			ResultBeforeAiFilter: map[string]interface{}{},
			Result:               map[string]interface{}{},
		},
	}

	result := convertNetlogs(logs)

	assert.Equal(t, len(logs), len(result))

	// Проверяем первый элемент
	assert.Equal(t, logs[0].ID, result[0].Id)
	assert.Equal(t, logs[0].ClientID, result[0].ClientId)
	assert.Equal(t, logs[0].AppName, result[0].AppName)
	assert.Equal(t, logs[0].Keywords, result[0].Keywords)
	assert.Equal(t, logs[0].Parameters, result[0].Parameters.AsMap())
	assert.Equal(t, *logs[0].Error, *result[0].Error)
	assert.Equal(t, int32(logs[0].NumBeforeAiFilter), result[0].NumBeforeAiFilter)
	assert.Equal(t, int32(logs[0].NumAfterAiFilter), result[0].NumAfterAiFilter)
	assert.Equal(t, logs[0].ResultBeforeAiFilter, result[0].ResultBeforeAiFilter.AsMap())
	assert.Equal(t, logs[0].Result, result[0].Result.AsMap())

	// Проверяем второй элемент (с nil error)
	assert.Equal(t, logs[1].ID, result[1].Id)
	assert.Nil(t, result[1].Error)
}

// Вспомогательные функции для создания указателей
func ptrString(s string) *string {
	return &s
}

func ptrInt64(i int64) *int64 {
	return &i
}

func ptrInt32(i int32) *int32 {
	return &i
}

func ptrInt(i int) *int {
	return &i
}

func ptrBool(b bool) *bool {
	return &b
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
