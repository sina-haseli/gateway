// internal/application/application.go
package application

import (
	"context"
	"fmt"
	"gateway/internal/domain"
	_ "gateway/internal/domain"
	"gateway/internal/infrastructure"
	pb "gateway/pkg/proto/protobuf"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Service represents the interface for user-related business logic.
type Service interface {
	GetRequest(ctx context.Context, request *http.Request) (*http.Response, error)
	GetGRPCRequest(ctx context.Context, request *http.Request) (*http.Response, error)
}

// ServiceImpl represents the concrete implementation of Service.
type ServiceImpl struct {
	loadBalancer domain.LoadBalancer
	infra        *infrastructure.Infrastructure
}

// NewService creates a new instance of Service.
func NewService(
	infra *infrastructure.Infrastructure,
	failureThreshold,
	successThreshold int,
	timeout time.Duration,
	svcs infrastructure.ServicesClient,
) Service {
	ctx := context.Background()
	// Retrieve services from the client
	services, err := svcs.Services(ctx)
	if err != nil {
		// Handle error gracefully
		panic(err)
	}
	// Initialize servers with circuit breakers
	var servers []domain.Server
	for _, service := range services {
		httpSender := domain.NewHttpServerSenderLogger(domain.NewHttpServerSender(service.Host, service.Port))
		cb := domain.NewCircuitBreaker(failureThreshold, successThreshold, timeout, httpSender)
		servers = append(servers, cb)
	}
	return &ServiceImpl{
		loadBalancer: domain.NewLoadBalancerLogger(domain.NewLoadBalancer(servers)),
		infra:        infra,
	}
}

func (s *ServiceImpl) GetRequest(ctx context.Context, request *http.Request) (*http.Response, error) {

	server, err := s.loadBalancer.Select(ctx)
	if err != nil {
		return nil, err
	}

	return server.SendRequest(ctx, request)
}

func (s *ServiceImpl) GetGRPCRequest(ctx context.Context, request *http.Request) (*http.Response, error) {
	grpcResponse, err := s.infra.GRPCClient.Call(ctx, "MyMethod", &pb.MyRequest{})
	if err != nil {
		return nil, err
	}

	fmt.Println("grpcResponse.(*pb.MyResponse).Result", grpcResponse.(*pb.MyResponse).GetResult())

	body := []byte(grpcResponse.(*pb.MyResponse).GetResult())
	httpResponse := &http.Response{
		StatusCode:    http.StatusOK,
		Proto:         request.Proto,
		ProtoMajor:    request.ProtoMajor,
		ProtoMinor:    request.ProtoMinor,
		Header:        make(http.Header),
		Body:          ioutil.NopCloser(strings.NewReader(string(body))),
		ContentLength: int64(len(body)),
		Request:       request,
	}

	httpResponse.Header.Set("Content-Type", "application/json")

	return httpResponse, nil
}
