package models

import "time"

type (
	Product struct {
		Id          int       `json:"id"`
		Name        string    `form:"name" json:"name" binding:"-" column:"name"`
		Description string    `json:"description" binding:"-" column:"description"`
		Price       int       `json:"price" binding:"-" column:"price"`
		Quantity    int       `json:"quantity" binding:"-" column:"quantity"`
		UserId      int       `json:"-"`
		CreatedAt   time.Time `json:"created_at" binding:"-"`
		UpdatedAt   time.Time `json:"updated_at" binding:"-"`
		Pagination
	}

	Products []Product
)

func (m *Product) ToLocal() {
	m.CreatedAt = m.CreatedAt.Local()
	m.UpdatedAt = m.UpdatedAt.Local()
}
