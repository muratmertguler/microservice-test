package product

import (
	"context"
)

type GetProductRequest struct {
}

type GetProductResponse struct {
}

type GetProductHandler struct {
}

func (h *GetProductHandler) handle(ctx context.Context, req *GetProductRequest) (*GetProductResponse, error) {
	return &GetProductResponse{}, nil
}
