package models

import "time"

type (
	User struct {
		Id        int       `json:"id" binding:"-"`
		FullName  string    `json:"full_name" binding:"-"`
		Email     string    `json:"email" binding:"-"`
		Password  string    `json:"password" binding:"-"`
		IsLogin   bool      `json:"is_login"`
		CreatedAt time.Time `json:"created_at" binding:"-"`
		UpdatedAt time.Time `json:"updated_at" binding:"-"`
	}

	ResetPassword struct {
		OldPassword string `json:"old_password" binding:""`
		NewPassword string `json:"new_password" binding:"-"`
	}
)
