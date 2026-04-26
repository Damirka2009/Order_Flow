package inventory

import (
	"context"
	"sync"
)

type Service struct {
	mu    sync.Mutex
	stock map[string]int32
}

func New() *Service {
	return &Service{
		stock: map[string]int32{
			"iphone": 10,
			"book":   0,
		},
	}
}

func (s *Service) CheckStock(ctx context.Context, item string, qty int32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	available, ok := s.stock[item]
	if !ok {
		return false
	}

	return available >= qty
}

func (s *Service) DecreaseStock(ctx context.Context, item string, qty int32) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	available, ok := s.stock[item]
	if !ok || available < qty {
		return false
	}

	s.stock[item] -= qty
	return true
}
