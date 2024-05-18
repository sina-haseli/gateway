package main

import (
	"gateway/pkg/mocks"
	pb "gateway/pkg/proto/protobuf"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serviceTwoMock := new(mocks.ServiceTwoMock)
	setupMockServiceTwo(serviceTwoMock)

	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, serviceTwoMock)

	log.Printf("Service Two Mock running on port %d", 50051)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func setupMockServiceTwo(src *mocks.ServiceTwoMock) {
	src.On("MyMethod", mock.Anything, &pb.MyRequest{}).Return(&pb.MyResponse{Result: "Mocked Response"}, nil)
}
