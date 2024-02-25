package models

type (
	Pagination struct {
		Page    int `form:"page" json:"-"`
		PerPage int `form:"per_page" json:"-"`
	}
)
