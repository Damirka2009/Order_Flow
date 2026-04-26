package service

import (
	"context"
	"errors"
	"fmt"
	"master/internal/domain"
	"master/internal/repository"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

type OrderService struct {
	repo *repository.OrdersRepository
}

func New(repo *repository.OrdersRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(ctx context.Context, item, category, currency string, price int64, quantity int32) *domain.Order {
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
	s.repo.Save(ctx, &newOrder)
	return &newOrder
}

func (s *OrderService) Get(ctx context.Context, id string) (*domain.Order, error) {
	order, ok := s.repo.Get(ctx, id)
	if !ok {
		return &domain.Order{}, status.Error(codes.NotFound, "order not found")
	}
	return order, nil
}

func (s *OrderService) Update(ctx context.Context, id string, incoming *domain.Order, mask *fieldmaskpb.FieldMask) (*domain.Order, error) {
	existing, ok := s.repo.Get(ctx, id)
	if !ok {
		return &domain.Order{}, errors.New("order not found")
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
	s.repo.Delete(ctx, id)
}

func (s *OrderService) List(ctx context.Context) []*domain.Order {
	return s.repo.List(ctx)
}
