package infrastructure

import (
	"context"
	"fmt"
	"gateway/config"
	"gateway/internal/domain"
	pb "gateway/pkg/proto/protobuf"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

type HttpSender interface {
	SendRequest(ctx context.Context, req *http.Request) (*http.Response, error)
}

type Infrastructure struct {
	GRPCClient    GRPCClient
	ServiceClient ServicesClient
}

func NewInfrastructure(gc GRPCClient, services ServicesClient) *Infrastructure {
	return &Infrastructure{
		GRPCClient:    gc,
		ServiceClient: services,
	}

}

type grpcClientImpl struct {
	host string
	port string
	// You can include any dependencies here.
}

func (g *grpcClientImpl) Call(ctx context.Context, method string, request interface{}) (interface{}, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", g.host, g.port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	// Create a client stub.
	client := pb.NewMyServiceClient(conn)

	// Make the gRPC call.
	response, err := client.MyMethod(ctx, &pb.MyRequest{})
	if err != nil {
		log.Fatalf("Error calling YourMethod: %v", err)
	}

	return response, nil
}

func NewGRPCClient(host string, port string) GRPCClient {
	// Initialize any dependencies here.
	return &grpcClientImpl{
		host: host,
		port: port,
	}
}

type GRPCClient interface {
	// Call sends a request to the gRPC server.
	Call(ctx context.Context, method string, request interface{}) (interface{}, error)
}

type ServicesClientImpl struct {
	// You can include any dependencies here.
	services []domain.Service
}

func (s ServicesClientImpl) Services(ctx context.Context) ([]domain.Service, error) {
	return s.services, nil
}

type ServicesClient interface {
	// Call sends a request to the gRPC server.
	Services(ctx context.Context) ([]domain.Service, error)
}

func NewServicesClient(services []domain.Service) ServicesClient {
	return &ServicesClientImpl{
		services: services,
	}
}

func DomainsToService(values config.Domain, address string) []domain.Service {

	var domains []string
	for _, value := range values.Main {
		domains = append(domains, value)
	}
	var services []domain.Service
	for _, d := range domains {
		services = append(services, domain.Service{
			Host: address,
			Port: d,
		})
	}
	return services
}

func HttpRequestToPbMyRequest(req *http.Request) *pb.MyRequest {
	return &pb.MyRequest{
		Data: req.Host,
	}
}
