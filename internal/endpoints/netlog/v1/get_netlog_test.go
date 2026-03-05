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

func TestController_GetNetlog(t *testing.T) {
	now := time.Now().UTC()
	errMsg := "error message"
	emptyErr := ""

	tests := []struct {
		name          string
		req           *pb.GetNetlogRequest
		mockSetup     func(mockService *mocks.MockService)
		expectedResp  *pb.GetNetlogResponse
		expectedError error
	}{
		{
			name: "successful get with all fields",
			req: &pb.GetNetlogRequest{
				Id: 123,
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetByID", mock.Anything, int64(123)).
					Return(&datastruct.Netlog{
						ID:        123,
						CreatedAt: now,
						ClientID:  "client-123",
						AppName:   "test-app",
						Keywords:  []string{"kw1", "kw2"},
						Parameters: map[string]interface{}{
							"param1": "value1",
							"param2": float64(42),
						},
						Error:             &errMsg,
						NumBeforeAiFilter: 100,
						NumAfterAiFilter:  50,
						ResultBeforeAiFilter: map[string]interface{}{
							"before1": "before value",
						},
						Result: map[string]interface{}{
							"after1": true,
						},
					}, nil)
			},
			expectedResp: &pb.GetNetlogResponse{
				Netlog: &pb.Netlog{
					Id:        123,
					CreatedAt: timestamppb.New(now),
					ClientId:  "client-123",
					AppName:   "test-app",
					Keywords:  []string{"kw1", "kw2"},
					Parameters: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"param1": structpb.NewStringValue("value1"),
							"param2": structpb.NewNumberValue(42),
						},
					},
					Error:             &errMsg,
					NumBeforeAiFilter: 100,
					NumAfterAiFilter:  50,
					ResultBeforeAiFilter: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"before1": structpb.NewStringValue("before value"),
						},
					},
					Result: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"after1": structpb.NewBoolValue(true),
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "successful get with nil error",
			req: &pb.GetNetlogRequest{
				Id: 456,
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetByID", mock.Anything, int64(456)).
					Return(&datastruct.Netlog{
						ID:        456,
						CreatedAt: now,
						ClientID:  "client-456",
						AppName:   "test-app-2",
						Keywords:  []string{"kw3"},
						Parameters: map[string]interface{}{
							"param3": "value3",
						},
						Error:                nil,
						NumBeforeAiFilter:    200,
						NumAfterAiFilter:     100,
						ResultBeforeAiFilter: map[string]interface{}{},
						Result: map[string]interface{}{
							"result": "success",
						},
					}, nil)
			},
			expectedResp: &pb.GetNetlogResponse{
				Netlog: &pb.Netlog{
					Id:        456,
					CreatedAt: timestamppb.New(now),
					ClientId:  "client-456",
					AppName:   "test-app-2",
					Keywords:  []string{"kw3"},
					Parameters: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"param3": structpb.NewStringValue("value3"),
						},
					},
					Error:                nil,
					NumBeforeAiFilter:    200,
					NumAfterAiFilter:     100,
					ResultBeforeAiFilter: &structpb.Struct{},
					Result: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"result": structpb.NewStringValue("success"),
						},
					},
				},
			},
			expectedError: nil,
		},
		{
			name: "successful get with empty error string",
			req: &pb.GetNetlogRequest{
				Id: 789,
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetByID", mock.Anything, int64(789)).
					Return(&datastruct.Netlog{
						ID:                   789,
						CreatedAt:            now,
						ClientID:             "client-789",
						AppName:              "test-app-3",
						Keywords:             []string{},
						Parameters:           map[string]interface{}{},
						Error:                &emptyErr,
						NumBeforeAiFilter:    0,
						NumAfterAiFilter:     0,
						ResultBeforeAiFilter: map[string]interface{}{},
						Result:               map[string]interface{}{},
					}, nil)
			},
			expectedResp: &pb.GetNetlogResponse{
				Netlog: &pb.Netlog{
					Id:                   789,
					CreatedAt:            timestamppb.New(now),
					ClientId:             "client-789",
					AppName:              "test-app-3",
					Keywords:             []string{},
					Parameters:           &structpb.Struct{},
					Error:                &emptyErr,
					NumBeforeAiFilter:    0,
					NumAfterAiFilter:     0,
					ResultBeforeAiFilter: &structpb.Struct{},
					Result:               &structpb.Struct{},
				},
			},
			expectedError: nil,
		},
		{
			name: "successful get with minimal fields",
			req: &pb.GetNetlogRequest{
				Id: 101112,
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetByID", mock.Anything, int64(101112)).
					Return(&datastruct.Netlog{
						ID:                   101112,
						CreatedAt:            now,
						ClientID:             "client-minimal",
						AppName:              "minimal-app",
						Keywords:             []string{},
						Parameters:           map[string]interface{}{},
						Error:                nil,
						NumBeforeAiFilter:    0,
						NumAfterAiFilter:     0,
						ResultBeforeAiFilter: map[string]interface{}{},
						Result:               map[string]interface{}{},
					}, nil)
			},
			expectedResp: &pb.GetNetlogResponse{
				Netlog: &pb.Netlog{
					Id:                   101112,
					CreatedAt:            timestamppb.New(now),
					ClientId:             "client-minimal",
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
			expectedError: nil,
		},
		{
			name: "not found - service returns nil, nil",
			req: &pb.GetNetlogRequest{
				Id: 999,
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetByID", mock.Anything, int64(999)).
					Return(nil, nil)
			},
			expectedResp:  nil,
			expectedError: datastruct.NewNotFound("netlog not found"),
		},
		{
			name: "service error - database error with AppError",
			req: &pb.GetNetlogRequest{
				Id: 777,
			},
			mockSetup: func(mockService *mocks.MockService) {
				dbErr := datastruct.NewInternalError("database connection failed")
				mockService.On("GetByID", mock.Anything, int64(777)).
					Return(nil, dbErr)
			},
			expectedResp:  nil,
			expectedError: datastruct.NewInternalError("database connection failed"),
		},
		{
			name: "service error - plain error",
			req: &pb.GetNetlogRequest{
				Id: 888,
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("GetByID", mock.Anything, int64(888)).
					Return(nil, errors.New("plain error"))
			},
			expectedResp:  nil,
			expectedError: errors.New("plain error"),
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
			resp, err := controller.GetNetlog(context.Background(), tt.req)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, resp)

				// Для not found проверяем что это действительно AppError с кодом NotFound
				if tt.name == "not found - service returns nil, nil" {
					var appErr *datastruct.AppError
					assert.True(t, errors.As(err, &appErr))
					assert.Equal(t, datastruct.CodeNotFound, appErr.Code)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedResp.GetNetlog().GetId(), resp.GetNetlog().GetId())
				assert.Equal(t, tt.expectedResp.GetNetlog().GetClientId(), resp.GetNetlog().GetClientId())
				assert.Equal(t, tt.expectedResp.GetNetlog().GetAppName(), resp.GetNetlog().GetAppName())
				assert.Equal(t, tt.expectedResp.GetNetlog().GetKeywords(), resp.GetNetlog().GetKeywords())
				assert.Equal(t, tt.expectedResp.GetNetlog().GetError(), resp.GetNetlog().GetError())
				assert.Equal(t, tt.expectedResp.GetNetlog().GetNumBeforeAiFilter(), resp.GetNetlog().GetNumBeforeAiFilter())
				assert.Equal(t, tt.expectedResp.GetNetlog().GetNumAfterAiFilter(), resp.GetNetlog().GetNumAfterAiFilter())

				// Проверяем Parameters
				if tt.expectedResp.GetNetlog().GetParameters() != nil {
					assert.Equal(t, tt.expectedResp.GetNetlog().GetParameters().AsMap(), resp.GetNetlog().GetParameters().AsMap())
				}

				// Проверяем ResultBeforeAiFilter
				if tt.expectedResp.GetNetlog().GetResultBeforeAiFilter() != nil {
					assert.Equal(t, tt.expectedResp.GetNetlog().GetResultBeforeAiFilter().AsMap(), resp.GetNetlog().GetResultBeforeAiFilter().AsMap())
				}

				// Проверяем Result
				if tt.expectedResp.GetNetlog().GetResult() != nil {
					assert.Equal(t, tt.expectedResp.GetNetlog().GetResult().AsMap(), resp.GetNetlog().GetResult().AsMap())
				}
			}

			// Verify that all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

// Тест для проверки что convertNetlog корректно работает
func TestConvertNetlog(t *testing.T) {
	now := time.Now().UTC()
	errMsg := "test error"

	dsNetlog := &datastruct.Netlog{
		ID:        123,
		CreatedAt: now,
		ClientID:  "client-123",
		AppName:   "test-app",
		Keywords:  []string{"kw1", "kw2"},
		Parameters: map[string]interface{}{
			"param1": "value1",
			"param2": float64(42),
		},
		Error:             &errMsg,
		NumBeforeAiFilter: 100,
		NumAfterAiFilter:  50,
		ResultBeforeAiFilter: map[string]interface{}{
			"before1": "before value",
		},
		Result: map[string]interface{}{
			"after1": true,
		},
	}

	pbNetlog := convertNetlog(dsNetlog)

	assert.Equal(t, dsNetlog.ID, pbNetlog.Id)
	assert.Equal(t, dsNetlog.CreatedAt.Unix(), pbNetlog.CreatedAt.AsTime().Unix())
	assert.Equal(t, dsNetlog.ClientID, pbNetlog.ClientId)
	assert.Equal(t, dsNetlog.AppName, pbNetlog.AppName)
	assert.Equal(t, dsNetlog.Keywords, pbNetlog.Keywords)
	assert.Equal(t, dsNetlog.Parameters, pbNetlog.Parameters.AsMap())
	assert.Equal(t, *dsNetlog.Error, *pbNetlog.Error)
	assert.Equal(t, int32(dsNetlog.NumBeforeAiFilter), pbNetlog.NumBeforeAiFilter)
	assert.Equal(t, int32(dsNetlog.NumAfterAiFilter), pbNetlog.NumAfterAiFilter)
	assert.Equal(t, dsNetlog.ResultBeforeAiFilter, pbNetlog.ResultBeforeAiFilter.AsMap())
	assert.Equal(t, dsNetlog.Result, pbNetlog.Result.AsMap())
}
