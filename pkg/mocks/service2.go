package mocks

import (
	"context"
	pb "gateway/pkg/proto/protobuf"
	"github.com/stretchr/testify/mock"
)

type ServiceTwoMock struct {
	mock.Mock
	pb.UnimplementedMyServiceServer
}

func (m *ServiceTwoMock) MyMethod(ctx context.Context, req *pb.MyRequest) (*pb.MyResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*pb.MyResponse), args.Error(1)
}
