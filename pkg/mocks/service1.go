package mocks

import "github.com/stretchr/testify/mock"
import "net/http"

// ServiceOneMock is a mock implementation of an HTTP service.
type ServiceOneMock struct {
	mock.Mock
}

// ServeHTTP mocks the HTTP response.
func (m *ServiceOneMock) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.Called(w, r)
}
