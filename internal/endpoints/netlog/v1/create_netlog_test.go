package netlog

import (
	"context"
	"testing"

	pb "github.com/oulabla/ai_app_netlog/gen/go/netlog/v1"
	"github.com/oulabla/ai_app_netlog/internal/datastruct"
	"github.com/oulabla/ai_app_netlog/internal/endpoints/netlog/v1/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestController_CreateNetlog(t *testing.T) {
	errMsg := "error message"
	now := timestamppb.Now()

	tests := []struct {
		name          string
		req           *pb.CreateNetlogRequest
		mockSetup     func(mockService *mocks.MockService)
		expectedID    int64
		expectedError error
	}{
		{
			name: "successful creation with all fields",
			req: &pb.CreateNetlogRequest{
				Netlog: &pb.Netlog{
					ClientId:  "test-client",
					AppName:   "test-app",
					Keywords:  []string{"keyword1", "keyword2"},
					CreatedAt: now,
					Parameters: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"param1": structpb.NewStringValue("value1"),
						},
					},
					Error:             &errMsg,
					NumBeforeAiFilter: 10,
					NumAfterAiFilter:  5,
					ResultBeforeAiFilter: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"before": structpb.NewStringValue("before-result"),
						},
					},
					Result: &structpb.Struct{
						Fields: map[string]*structpb.Value{
							"after": structpb.NewStringValue("after-result"),
						},
					},
				},
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("Create", mock.Anything, mock.MatchedBy(func(netlog *datastruct.Netlog) bool {
					return netlog.ClientID == "test-client" &&
						netlog.AppName == "test-app" &&
						len(netlog.Keywords) == 2 &&
						netlog.Keywords[0] == "keyword1" &&
						netlog.Keywords[1] == "keyword2" &&
						netlog.Parameters["param1"] == "value1" &&
						*netlog.Error == "error message" &&
						netlog.NumBeforeAiFilter == 10 &&
						netlog.NumAfterAiFilter == 5 &&
						netlog.ResultBeforeAiFilter["before"] == "before-result" &&
						netlog.Result["after"] == "after-result"
				})).Return(int64(123), nil)
			},
			expectedID:    123,
			expectedError: nil,
		},
		{
			name: "successful creation without optional fields",
			req: &pb.CreateNetlogRequest{
				Netlog: &pb.Netlog{
					ClientId: "test-client",
					AppName:  "test-app",
					Result:   &structpb.Struct{},
				},
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("Create", mock.Anything, mock.MatchedBy(func(netlog *datastruct.Netlog) bool {
					return netlog.ClientID == "test-client" &&
						netlog.AppName == "test-app" &&
						len(netlog.Keywords) == 0 &&
						len(netlog.Parameters) == 0 &&
						*netlog.Error == "" && // пустая строка по умолчанию
						netlog.NumBeforeAiFilter == 0 &&
						netlog.NumAfterAiFilter == 0 &&
						len(netlog.ResultBeforeAiFilter) == 0 &&
						len(netlog.Result) == 0
				})).Return(int64(456), nil)
			},
			expectedID:    456,
			expectedError: nil,
		},
		{
			name: "successful creation with nil error",
			req: &pb.CreateNetlogRequest{
				Netlog: &pb.Netlog{
					ClientId: "test-client",
					AppName:  "test-app",
					Error:    nil,
					Result:   &structpb.Struct{},
				},
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("Create", mock.Anything, mock.MatchedBy(func(netlog *datastruct.Netlog) bool {
					return netlog.ClientID == "test-client" &&
						netlog.AppName == "test-app" &&
						*netlog.Error == "" // пустая строка когда Error = nil
				})).Return(int64(789), nil)
			},
			expectedID:    789,
			expectedError: nil,
		},
		{
			name: "service error",
			req: &pb.CreateNetlogRequest{
				Netlog: &pb.Netlog{
					ClientId: "test-client",
					AppName:  "test-app",
					Result:   &structpb.Struct{},
				},
			},
			mockSetup: func(mockService *mocks.MockService) {
				mockService.On("Create", mock.Anything, mock.Anything).Return(int64(0), assert.AnError)
			},
			expectedID:    0,
			expectedError: assert.AnError,
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
			resp, err := controller.CreateNetlog(context.Background(), tt.req)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.expectedID, resp.Id)
			}

			// Verify that all mock expectations were met
			mockService.AssertExpectations(t)
		})
	}
}

func TestController_CreateNetlog_FieldMapping(t *testing.T) {
	errMsg := "error message"
	mockService := new(mocks.MockService)

	// Create request with all fields populated
	req := &pb.CreateNetlogRequest{
		Netlog: &pb.Netlog{
			ClientId: "client-123",
			AppName:  "app-name",
			Keywords: []string{"kw1", "kw2"},
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
	}

	// Setup mock to capture the datastruct
	mockService.On("Create", mock.Anything, mock.MatchedBy(func(netlog *datastruct.Netlog) bool {
		// Проверяем все поля
		return netlog.ClientID == "client-123" &&
			netlog.AppName == "app-name" &&
			len(netlog.Keywords) == 2 &&
			netlog.Keywords[0] == "kw1" &&
			netlog.Keywords[1] == "kw2" &&
			netlog.Parameters["param1"] == "value1" &&
			netlog.Parameters["param2"] == float64(42) &&
			*netlog.Error == "error message" &&
			netlog.NumBeforeAiFilter == 100 &&
			netlog.NumAfterAiFilter == 50 &&
			netlog.ResultBeforeAiFilter["before1"] == "before value" &&
			netlog.Result["after1"] == true
	})).Return(int64(123), nil)

	controller := &Controller{
		service: mockService,
	}

	resp, err := controller.CreateNetlog(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), resp.Id)
	mockService.AssertExpectations(t)
}

func TestController_CreateNetlog_ErrorHandling(t *testing.T) {
	mockService := new(mocks.MockService)

	// Тест с nil Error
	reqWithNilError := &pb.CreateNetlogRequest{
		Netlog: &pb.Netlog{
			ClientId: "client-123",
			AppName:  "app-name",
			Error:    nil,
			Result:   &structpb.Struct{},
		},
	}

	mockService.On("Create", mock.Anything, mock.MatchedBy(func(netlog *datastruct.Netlog) bool {
		// Проверяем что Error не nil и содержит пустую строку
		return netlog.Error != nil && *netlog.Error == ""
	})).Return(int64(123), nil)

	controller := &Controller{
		service: mockService,
	}

	resp, err := controller.CreateNetlog(context.Background(), reqWithNilError)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), resp.Id)
	mockService.AssertExpectations(t)
}
