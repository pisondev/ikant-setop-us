package fish

import "time"

type FishType struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ImageURL    *string   `json:"image_url"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInput struct {
	Name        string
	ImageURL    *string
	Description *string
}
