package tests

import (
	_ "context"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	_ "google.golang.org/grpc/resolver"

	"gateway/pkg/mocks"
)

func TestGateway(t *testing.T) {
	// Mock HTTP services for Service One
	service1Mock1 := new(mocks.ServiceOneMock)
	service1Mock2 := new(mocks.ServiceOneMock)
	service1Mock3 := new(mocks.ServiceOneMock)

	service1Mock1.On("ServeHTTP", mock.Anything, mock.Anything).Return(http.StatusOK, "Service One Response 1")
	service1Mock2.On("ServeHTTP", mock.Anything, mock.Anything).Return(http.StatusOK, "Service One Response 2")
	service1Mock3.On("ServeHTTP", mock.Anything, mock.Anything).Return(http.StatusOK, "Service One Response 3")

	// Mock gRPC service for Service Two
	// Start mock services
	setupMockServiceOne(service1Mock1, 8081)
	setupMockServiceOne(service1Mock2, 8082)
	setupMockServiceOne(service1Mock3, 8083)

	time.Sleep(1 * time.Second) // Wait for servers to start

	// Create http client to test Service One
	RS1, err := http.Get("http://localhost:8081")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, RS1.StatusCode)

	RS2, err := http.Get("http://localhost:8082")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, RS2.StatusCode)

	RS3, err := http.Get("http://localhost:8083")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, RS3.StatusCode)

}

func setupMockServiceOne(src *mocks.ServiceOneMock, port int) {
	src.On("ServeHTTP", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		w := args.Get(0).(http.ResponseWriter)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service One Response" + " " + fmt.Sprint(port)))
	})

	mux := http.NewServeMux()
	mux.Handle("/", src)

	go func() {
		log.Printf("Service One Mock running on port %d", port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
			log.Fatalf("Failed to start Service One Mock on port %d: %v", port, err)
		}
	}()
}

