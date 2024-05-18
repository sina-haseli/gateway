package pkg

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "gateway/pkg/proto/protobuf"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedMyServiceServer
}

func (s *Server) MyMethod(ctx context.Context, req *pb.MyRequest) (*pb.MyResponse, error) {
	// Business logic
	return &pb.MyResponse{Result: "Hello " + req.Data}, nil
}

func StartGRPCServer(host, port string) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterMyServiceServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
