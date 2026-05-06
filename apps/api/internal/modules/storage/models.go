package storage

import "time"

type ColdStorage struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	LocationLabel *string   `json:"location_label"`
	Description   *string   `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateInput struct {
	Name          string
	LocationLabel *string
	Description   *string
}
