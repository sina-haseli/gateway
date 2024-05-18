package http

import (
	"gateway/internal/application"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Handler represents the HTTP handler for user-related endpoints.
type Handler interface {
	GetRequest(c echo.Context) error
	GetGRPCRequest(c echo.Context) error
}

// HandlerImpl represents the concrete implementation of Handler.
type HandlerImpl struct {
	handler application.Service
}

// NewHandler creates a new instance of Handler.
func NewHandler(serv application.Service) Handler {
	return &HandlerImpl{handler: serv}
}

// GetUser handles the HTTP request to get a user by ID.
func (h *HandlerImpl) GetRequest(c echo.Context) error {
	_, err := h.handler.GetRequest(c.Request().Context(), c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get Request"})
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *HandlerImpl) GetGRPCRequest(c echo.Context) error {
	res, err := h.handler.GetGRPCRequest(c.Request().Context(), c.Request())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get Request"})
	}

	return c.JSON(http.StatusOK, res)
}
