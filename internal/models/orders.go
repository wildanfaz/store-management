package models

import "time"

type (
	Order struct {
		Id        int       `json:"id"`
		ProductId int       `form:"product_id" json:"product_id,omitempty" binding:"-"`
		UserId    int       `json:"user_id,omitempty"`
		Quantity  int       `json:"quantity" binding:"-"`
		Status    string    `form:"status" json:"status"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Pagination
	}

	OrderResponse struct {
		Order   Order   `json:"order"`
		Product Product `json:"product"`
	}

	Orders []OrderResponse
)

func (m *OrderResponse) ToLocal() {
	m.Order.CreatedAt = m.Order.CreatedAt.Local()
	m.Order.UpdatedAt = m.Order.UpdatedAt.Local()

	m.Product.CreatedAt = m.Product.CreatedAt.Local()
	m.Product.UpdatedAt = m.Product.UpdatedAt.Local()
}
