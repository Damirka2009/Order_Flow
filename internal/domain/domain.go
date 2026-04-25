package domain

type Order struct {
	Id       string
	Item     string
	Category string
	Currency string
	Price    int64
	Quantity int32
	Is_stock bool
}
