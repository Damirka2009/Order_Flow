package service

import (
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

func (s *OrderService) Create(item, category, currency string, price int64, quantity int32) *domain.Order {
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
	s.repo.Save(&newOrder)
	return &newOrder
}

func (s *OrderService) Get(id string) (*domain.Order, error) {
	order, ok := s.repo.Get(id)
	if !ok {
		return &domain.Order{}, status.Error(codes.NotFound, "order not found")
	}
	return order, nil
}

func (s *OrderService) Update(id string, incoming *domain.Order, mask *fieldmaskpb.FieldMask) (*domain.Order, error) {
	existing, ok := s.repo.Get(id)
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
	s.repo.Update(existing)
	return existing, nil
}

func (s *OrderService) Delete(id string) {
	s.repo.Delete(id)
}

func (s *OrderService) List() []*domain.Order {
	return s.repo.List()
}
