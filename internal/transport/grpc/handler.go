package grpc

import (
	"context"

	"master/internal/domain"
	"master/internal/service"
	api "master/pkg/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	api.UnimplementedOrderServiceServer
	svc *service.OrderService
}

func New(svc *service.OrderService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateOrder(ctx context.Context, req *api.CreateOrderRequest) (*api.OrderResponse, error) {
	order := h.svc.Create(
		ctx,
		req.Item,
		req.Category,
		req.Currency,
		req.Price,
		req.Quantity,
	)

	return &api.OrderResponse{
		Order: &api.Order{
			Id:       order.Id,
			Item:     order.Item,
			Category: order.Category,
			Currency: order.Currency,
			Price:    order.Price,
			Quantity: order.Quantity,
			IsStock:  order.Is_stock,
		},
	}, nil
}

func (h *Handler) GetOrder(ctx context.Context, req *api.GetOrderRequest) (*api.OrderResponse, error) {
	order, err := h.svc.Get(ctx, req.Id)
	if err != nil {
		return &api.OrderResponse{}, err
	}
	return &api.OrderResponse{
		Order: &api.Order{
			Id:       order.Id,
			Item:     order.Item,
			Category: order.Category,
			Currency: order.Currency,
			Price:    order.Price,
			Quantity: order.Quantity,
			IsStock:  order.Is_stock,
		},
	}, nil
}

func (h *Handler) UpdateOrder(ctx context.Context, req *api.UpdateOrderRequest) (*api.OrderResponse, error) {
	incoming := &domain.Order{
		Id:       req.Order.Id,
		Item:     req.Order.Item,
		Category: req.Order.Category,
		Currency: req.Order.Currency,
		Price:    req.Order.Price,
		Quantity: req.Order.Quantity,
		Is_stock: req.Order.IsStock,
	}
	order, err := h.svc.Update(ctx, req.Id, incoming, req.UpdateMask)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.OrderResponse{
		Order: &api.Order{
			Id:       order.Id,
			Item:     order.Item,
			Category: order.Category,
			Currency: order.Currency,
			Price:    order.Price,
			Quantity: order.Quantity,
			IsStock:  order.Is_stock,
		},
	}, nil
}

func (h *Handler) DeleteOrder(ctx context.Context, req *api.DeleteOrderRequest) (*api.Empty, error) {
	h.svc.Delete(ctx, req.Id)
	return &api.Empty{}, nil
}

func (h *Handler) OrdersList(ctx context.Context, _ *api.Empty) (*api.OrdersListResponse, error) {
	orders := h.svc.List(ctx)

	var pbOrders []*api.Order

	for _, o := range orders {
		pbOrders = append(pbOrders, &api.Order{
			Id:       o.Id,
			Item:     o.Item,
			Category: o.Category,
			Currency: o.Currency,
			Price:    o.Price,
			Quantity: o.Quantity,
			IsStock:  o.Is_stock,
		})
	}

	return &api.OrdersListResponse{
		Orders: pbOrders,
	}, nil
}
