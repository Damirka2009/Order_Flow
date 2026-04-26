package grpc

import (
	"context"

	"master/internal/inventory"
	api "master/pkg/api"

	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	api.UnimplementedInventoryServiceServer
	svc *inventory.Service
}

func New(svc *inventory.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CheckStock(ctx context.Context, req *api.CheckStockRequest) (*api.CheckStockResponse, error) {
	ok := h.svc.CheckStock(ctx, req.Item, req.Quantity)

	return &api.CheckStockResponse{
		Available: ok,
	}, nil
}

func (h *Handler) DecreaseStock(ctx context.Context, req *api.DecreaseStockRequest) (*emptypb.Empty, error) {
	h.svc.DecreaseStock(ctx, req.Item, req.Quantity)
	return &emptypb.Empty{}, nil
}
