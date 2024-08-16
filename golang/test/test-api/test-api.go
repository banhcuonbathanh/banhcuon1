package testapi

import (
	"context"
	"net/http"
	"net/http/httptest"

	"english-ai-full/ecomm-api/handler"
	"testing"

	pb "english-ai-full/ecomm-grpc/proto"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockEcommUserClient is a mock for the EcommUserClient
type MockEcommUserClient struct {
	mock.Mock
}

// Implement all methods of pb.EcommUserClient interface
func (m *MockEcommUserClient) CreateUser(ctx context.Context, in *pb.UserReq, opts ...grpc.CallOption) (*pb.UserReq, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.UserReq), args.Error(1)
}

func (m *MockEcommUserClient) SaveUser(ctx context.Context, in *pb.UserReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockEcommUserClient) UpdateUser(ctx context.Context, in *pb.UserReq, opts ...grpc.CallOption) (*pb.UserReq, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.UserReq), args.Error(1)
}

func (m *MockEcommUserClient) DeleteUser(ctx context.Context, in *pb.UserReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*emptypb.Empty), args.Error(1)
}

func (m *MockEcommUserClient) FindAllUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*pb.UserList, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.UserList), args.Error(1)
}

func (m *MockEcommUserClient) FindByEmail(ctx context.Context, in *pb.UserReq, opts ...grpc.CallOption) (*pb.UserReq, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.UserReq), args.Error(1)
}

func (m *MockEcommUserClient) FindUsersByPage(ctx context.Context, in *pb.PageRequest, opts ...grpc.CallOption) (*pb.UserList, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.UserList), args.Error(1)
}

func (m *MockEcommUserClient) Login(ctx context.Context, in *pb.LoginRequest, opts ...grpc.CallOption) (*pb.UserReq, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.UserReq), args.Error(1)
}

func (m *MockEcommUserClient) Register(ctx context.Context, in *pb.UserReq, opts ...grpc.CallOption) (*pb.RegisterResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.RegisterResponse), args.Error(1)
}

// Rest of your test functions remain the same
// ...

func TestHomeRoute(t *testing.T) {
	handlerTest := handler.NewHandler(&MockEcommUserClient{}, "test-secret")
	r := handler.RegisterRoutes(handlerTest)

	// Create a new HTTP request to the home route
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record the response
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Check the status code is what we expect
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect
	expected := "Server is running"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
