package domain

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"time"
)

// CircuitBreakerImpl represents a circuit breaker.
type CircuitBreakerImpl struct {
	mu               sync.Mutex
	state            string
	failureThreshold int
	successThreshold int
	timeout          time.Duration
	failureCount     int
	successCount     int
	lastFailureTime  time.Time
	server           Server
}

func (c *CircuitBreakerImpl) SendRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	var resp *http.Response
	fn := func() error {
		var err error
		resp, err = c.server.SendRequest(ctx, req)
		if err != nil {
			return err
		}
		return nil
	}

	err := c.execute(fn)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CircuitBreakerImpl) IsHealthy(ctx context.Context) bool {
	fn := func() error {
		if !c.server.IsHealthy(ctx) {
			return errors.New("server is unhealthy")
		}
		return nil
	}

	err := c.execute(fn)
	if err != nil {
		return false
	}

	return true
}

// NewCircuitBreaker creates a new circuit breaker with the given parameters.
func NewCircuitBreaker(failureThreshold, successThreshold int, timeout time.Duration, server Server) Server {
	return &CircuitBreakerImpl{
		state:            "CLOSED",
		failureThreshold: failureThreshold,
		successThreshold: successThreshold,
		timeout:          timeout * time.Second,
		server:           server,
	}
}

func (c *CircuitBreakerImpl) execute(fn func() error) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch c.state {
	case "OPEN":
		// If the circuit is open, check if it's time to transition to half-open state
		if time.Since(c.lastFailureTime) >= c.timeout {
			c.state = "HALF_OPEN"
		} else {
			return errors.New("circuit is open")
		}
	case "HALF_OPEN":
		// If the circuit is in half-open state, allow a single attempt
		err := fn()
		if err == nil {
			// If the attempt is successful, transition to closed state
			c.state = "CLOSED"
			c.successCount++
			return nil
		} else {
			// If the attempt fails, transition back to open state
			c.state = "OPEN"
			c.failureCount++
			c.lastFailureTime = time.Now()
			return err
		}
	case "CLOSED":
		// If the circuit is closed, execute the function
		err := fn()
		if err == nil {
			c.successCount++
			// If the function is successful, reset failure count
			if c.successCount >= c.successThreshold {
				c.failureCount = 0
			}
		} else {
			c.failureCount++
			c.lastFailureTime = time.Now()
			// If failure threshold is reached, transition to open state
			if c.failureCount >= c.failureThreshold {
				c.state = "OPEN"
			}
		}
		return err
	default:
		return errors.New("invalid state")
	}

	return nil
}
