package healthcheck

import (
	"context"
	_ "github.com/gofiber/fiber/v2/middleware/recover"
)

type healthCheckHandler struct {
}

type HealthCheckResponse struct {
	Status string `json:"status"`
}
type HealthCheckRequest struct {
}

func NewHealthCheckHandler() *healthCheckHandler {
	return &healthCheckHandler{}
}

func (h *healthCheckHandler) Handle(ctx context.Context, req *HealthCheckRequest) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{Status: "OK"}, nil
}
