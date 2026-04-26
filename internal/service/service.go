package service

import (
	"context"
	"errors"
	"fmt"
	"master/internal/domain"
	"master/internal/repository"
	"time"

	api "master/pkg/api"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type OrderService struct {
	repo            *repository.OrdersRepository
	inventoryClient api.InventoryServiceClient
}

func New(repo *repository.OrdersRepository, inv api.InventoryServiceClient) *OrderService {
	return &OrderService{
		repo:            repo,
		inventoryClient: inv,
	}
}

func (s *OrderService) Create(ctx context.Context, item, category, currency string, price int64, quantity int32) *domain.Order {
	resp, err := s.inventoryClient.CheckStock(ctx, &api.CheckStockRequest{
		Item:     item,
		Quantity: quantity,
	})

	if err != nil || !resp.Available {
		return &domain.Order{}
	}

	var is_stock bool
	if quantity > 0 {
		is_stock = true
	} else {
		is_stock = false
	}
	newOrder := domain.Order{
		Id:       uuid.New().String(),
		Item:     item,
		Category: category,
		Currency: currency,
		Price:    price,
		Quantity: quantity,
		Is_stock: is_stock,
	}
	err = withRetry(ctx, 2, func() error {
		return s.repo.Save(ctx, &newOrder)
	})

	_, _ = s.inventoryClient.DecreaseStock(ctx, &api.DecreaseStockRequest{
		Item:     item,
		Quantity: quantity,
	})

	if err != nil {
		return &domain.Order{}
	}
	return &newOrder
}

func (s *OrderService) Get(ctx context.Context, id string) (*domain.Order, error) {
	var order *domain.Order

	err := withRetry(ctx, 2, func() error {
		var innerOk bool
		order, innerOk = s.repo.Get(ctx, id)
		if !innerOk {
			return fmt.Errorf("not found")
		}
		return nil
	})

	if err != nil {
		return &domain.Order{}, status.Error(codes.NotFound, "fallback: order not available")
	}

	return order, nil
}

func (s *OrderService) Update(ctx context.Context, id string, incoming *domain.Order, mask *fieldmaskpb.FieldMask) (*domain.Order, error) {
	existing, ok := s.repo.Get(ctx, id)
	if !ok {
		return &domain.Order{}, fmt.Errorf("fallback: update failed")
	}

	if mask == nil || len(mask.Paths) == 0 {
		return &domain.Order{}, errors.New("update_mask is required")
	}

	for _, path := range mask.Paths {
		switch path {
		case "item":
			existing.Item = incoming.Item
		case "category":
			existing.Category = incoming.Category
		case "currency":
			existing.Currency = incoming.Currency
		case "price":
			existing.Price = incoming.Price
		case "quantity":
			existing.Quantity = incoming.Quantity
		case "is_stock":
			existing.Is_stock = incoming.Is_stock
		default:
			return nil, fmt.Errorf("unknown field: %s", path)
		}
	}
	s.repo.Update(ctx, existing)
	return existing, nil
}

func (s *OrderService) Delete(ctx context.Context, id string) {
	_ = withRetry(ctx, 2, func() error {
		return s.repo.Delete(ctx, id)
	})
}

func (s *OrderService) List(ctx context.Context) []*domain.Order {
	var orders []*domain.Order

	err := withRetry(ctx, 2, func() error {
		orders = s.repo.List(ctx)
		if orders == nil {
			return fmt.Errorf("failed")
		}
		return nil
	})

	if err != nil {
		return []*domain.Order{}
	}

	return orders
}

func withRetry(ctx context.Context, attempts int, fn func() error) error {
	var err error

	for i := 0; i < attempts; i++ {
		_, cancel := context.WithTimeout(ctx, 2*time.Second)

		err = fn()
		cancel()

		if err == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return err
}
