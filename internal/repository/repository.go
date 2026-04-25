package repository

import (
	"master/internal/domain"
	"sync"
)

type OrdersRepository struct {
	mu     sync.RWMutex
	orders map[string]*domain.Order
}

func New() *OrdersRepository {
	return &OrdersRepository{
		orders: make(map[string]*domain.Order),
	}
}

func (r *OrdersRepository) Save(order *domain.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.Id] = order
}

func (r *OrdersRepository) Get(id string) (*domain.Order, bool) {
	r.mu.RLock()
	defer r.mu.Unlock()
	order, ok := r.orders[id]
	return order, ok
}

func (r *OrdersRepository) Update(order *domain.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.Id] = order
}

func (r *OrdersRepository) Delete(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.orders, id)
}

func (r *OrdersRepository) List() []*domain.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()

	res := []*domain.Order{}
	for _, order := range r.orders {
		res = append(res, order)
	}
	return res
}
