// domain/loadbalancer.go
package domain

import (
	"context"
	"fmt"
	"sync"
)

type LoadBalancer interface {
	Select(ctx context.Context) (Server, error)
}

type loadBalancerImpl struct {
	mu      sync.Mutex
	servers []Server
	current int
}

func NewLoadBalancer(servers []Server) LoadBalancer {
	return &loadBalancerImpl{
		servers: servers,
		current: 0,
	}
}

func (lb *loadBalancerImpl) Select(ctx context.Context) (Server, error) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	if len(lb.servers) == 0 {
		return nil, fmt.Errorf("no server available")
	}

	healthyServers := make([]Server, 0, len(lb.servers))

	for _, server := range lb.servers {
		if server.IsHealthy(ctx) {
			healthyServers = append(healthyServers, server)
		}
	}

	if len(healthyServers) == 0 {
		return nil, fmt.Errorf("no healthy server available")
	}

	if lb.current >= len(healthyServers) {
		lb.current = 0
	}

	server := healthyServers[lb.current]
	lb.current = (lb.current + 1) % len(healthyServers)

	return server, nil
}

type loadBalancerImplLogger struct {
	base LoadBalancer
}

func NewLoadBalancerLogger(base LoadBalancer) LoadBalancer {
	return &loadBalancerImplLogger{
		base: base,
	}
}

func (l loadBalancerImplLogger) Select(ctx context.Context) (Server, error) {
	server, err := l.base.Select(ctx)
	if err != nil {
		fmt.Printf("LoadBalancer: Select error: %s\n", err)
	} else {
		fmt.Printf("LoadBalancer: Select: %s\n", server)
	}
	return server, err
}
