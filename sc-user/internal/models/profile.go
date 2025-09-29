package models

import "time"

type UserProfile struct {
	ID        string    `json:"id"`
	AuthID    string    `json:"auth_id"`
	FullName  string    `json:"full_name"`
	Bio       *string   `json:"bio,omitempty"`
	PhotoLink *string   `json:"photo_link,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateUserProfileInput struct {
	AuthID    string  `json:"auth_id"`
	FullName  string  `json:"full_name"`
	Bio       *string `json:"bio,omitempty"`
	PhotoLink *string `json:"photo_link,omitempty"`
}

type UpdateUserProfileInput struct {
	ID        string  `json:"id"`
	FullName  *string `json:"full_name,omitempty"`
	Bio       *string `json:"bio,omitempty"`
	PhotoLink *string `json:"photo_link,omitempty"`
}

type UpdateUserProfileParams struct {
	ID        string  `json:"id"`
	FullName  string  `json:"full_name"`
	Bio       *string `json:"bio,omitempty"`
	PhotoLink *string `json:"photo_link,omitempty"`
}
